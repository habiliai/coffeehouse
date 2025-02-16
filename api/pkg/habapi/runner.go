package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	haberrors "github.com/habiliai/habiliai/api/pkg/errors"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"slices"
	"time"
)

type runRequest struct {
	threadId uint
}

func (s *server) run(ctx context.Context, threadId uint) error {
	select {
	case <-ctx.Done():
		return nil
	case s.runReqCh <- runRequest{threadId: threadId}:
		logger.Info("send doRunner request", "thread_id", threadId)
	}

	return nil
}

func (s *server) doRunner(ctx context.Context, workerName string) {
	defer logger.Info("stop goroutine", "worker", workerName)
	logger.Info("start goroutine", "worker", workerName)

	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-s.runReqCh:
			{
				if !ok {
					return
				}
				logger.Debug("doRunner request", "req", req)
				if err := s.runnerMain(ctx, req); err != nil {
					logger.Warn("failed to doRunner step", "err", err)
				}
			}
		}
	}
}

func (s *server) runnerMain(ctx context.Context, req runRequest) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eventListener := helpers.NewEventListener()
	ctx = helpers.WithEventListener(ctx, eventListener)

	tx := s.db.WithContext(ctx)
	ctx = helpers.WithTx(ctx, tx)

	var thread domain.Thread
	if err := tx.
		Preload("Mission").
		Preload("Mission.Steps").
		Preload("Mission.Steps.Actions").
		Preload("Mission.Steps.Actions.Agent").
		First(&thread, req.threadId).Error; err != nil {
		return errors.Wrapf(err, "failed to find thread with id %d", req.threadId)
	}

	step, err := thread.GetCurrentStep()
	if err != nil {
		return err
	}

	lastMessage, err := s.getLastMessage(ctx, thread.OpenaiThreadId)
	if err != nil {
		return err
	}

	var agents []*domain.Agent
	if lastMessage != nil && lastMessage.Role == Message_USER && len(lastMessage.Mentions) > 0 {
		for _, mention := range lastMessage.Mentions {
			var agent domain.Agent
			if err := tx.First(&agent, "name = ?", mention).Error; err != nil {
				return err
			}
			agents = append(agents, &agent)
		}
	} else if lastMessage != nil && lastMessage.Role == Message_USER {
		return errors.Errorf("user message must mention at least an agent")
	} else if lastMessage == nil || lastMessage.Role == Message_ASSISTANT {
		for _, action := range step.Actions {
			agents = append(agents, &action.Agent)
		}
	} else {
		return errors.Errorf("last message role is invalid")
	}

	slices.SortFunc(agents, func(a, b *domain.Agent) int {
		if a.IncludeQuestionIntent && b.IncludeQuestionIntent {
			return 0
		} else if !a.IncludeQuestionIntent && !b.IncludeQuestionIntent {
			return 0
		} else if !a.IncludeQuestionIntent && b.IncludeQuestionIntent {
			return -1
		} else {
			return 1
		}
	})

	logger.Debug("step actions", "actions", step.Actions)

	var run *openai.Run
	for _, action := range step.Actions {
		run, err = s.openai.Beta.Threads.Runs.New(ctx, thread.OpenaiThreadId, openai.BetaThreadRunNewParams{
			AssistantID: openai.F(action.Agent.AssistantId),
		})
		if err != nil {
			return errors.Wrapf(err, "failed to doRunner thread")
		}

		var actionWork domain.ActionWork
		if err := tx.Preload("Action").First(&actionWork, "action_id = ? AND thread_id = ?", action.ID, thread.ID).Error; err != nil {
			return errors.Wrapf(err, "failed to find action work")
		}

		var agentWork domain.AgentWork
		if err := tx.Preload("Agent").First(&agentWork, "agent_id = ? AND thread_id = ?", action.Agent.ID, thread.ID).Error; err != nil {
			return errors.Wrapf(err, "failed to find agent work")
		}

		agentWork.Status = domain.AgentStatusWorking
		if err := agentWork.Save(tx); err != nil {
			return err
		}

		for interrupt := false; !interrupt; {
			logger.Info("polling doRunner status", "doRunner", run)
			run, err = s.openai.Beta.Threads.Runs.PollStatus(ctx, thread.OpenaiThreadId, run.ID, 500)
			if err != nil {
				errors.Wrapf(err, "failed to get doRunner")
			}

			thread.CurrentRunId = run.ID
			if err := thread.Save(tx); err != nil {
				return err
			}

			switch run.Status {
			case openai.RunStatusCompleted:
				logger.Debug("doRunner completed", "doRunner", run)
				if agentWork.Status != domain.AgentStatusIdle {
					if !agentWork.Agent.IncludeQuestionIntent {
						agentWork.Status = domain.AgentStatusIdle
					} else {
						agentWork.Status = domain.AgentStatusWaiting
					}
					if err := agentWork.Save(tx); err != nil {
						return err
					}
				}
				interrupt = true
			case openai.RunStatusFailed:
				return errors.Wrapf(haberrors.ErrRuntime, "doRunner failed: %s", run.LastError.Message)
			case openai.RunStatusIncomplete:
				return errors.Wrapf(haberrors.ErrRuntime, "Run incomplete: %s", run.IncompleteDetails.Reason)
			case openai.RunStatusExpired:
				return errors.Wrapf(haberrors.ErrRuntime, "Run expired. expires_at: %s", time.Unix(run.ExpiresAt, 0))
			case openai.RunStatusCancelled:
				return errors.Wrapf(haberrors.ErrRuntime, "Run cancelled")
			case openai.RunStatusRequiresAction:
				run, err = s.requiresCallback(ctx, run, &thread, &actionWork, &agentWork)
				if err != nil {
					return err
				}
			default:
				return errors.Wrapf(haberrors.ErrBadRequest, "invalid thread doRunner status: received %s", run.Status)
			}
		}
	}

	eventListener.Emit(ctx, helpers.EventTypeEndRunning)

	return nil
}
