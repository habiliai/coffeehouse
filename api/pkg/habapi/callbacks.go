package habapi

import (
	"context"
	"encoding/json"
	"github.com/habiliai/habiliai/api/pkg/callbacks"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/mokiat/gog"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

func (s *server) requiresCallback(ctx context.Context, run *openai.Run, thread *domain.Thread, actionWork *domain.ActionWork, agentWork *domain.AgentWork) (*openai.Run, error) {
	if len(run.RequiredAction.SubmitToolOutputs.ToolCalls) == 0 {
		return nil, errors.Errorf("no tool calls found")
	}

	outputs := make([]openai.BetaThreadRunSubmitToolOutputsParamsToolOutput, 0, len(run.RequiredAction.SubmitToolOutputs.ToolCalls))

	for _, toolCall := range run.RequiredAction.SubmitToolOutputs.ToolCalls {
		var (
			err    error
			output struct {
				Success bool   `json:"success"`
				Reason  string `json:"reason,omitempty"`
				Result  any    `json:"result.omitempty"`
			}
		)

		var outputErr error
		switch strings.ToLower(toolCall.Function.Name) {
		case "done_agent": // this special case is used to mark the agent as idle
			output.Result, outputErr = s.doneAgent(ctx, thread, agentWork, actionWork)
		default:
			metadata := callbacks.Metadata{
				AgentWork:  agentWork,
				ActionWork: actionWork,
			}
			output.Result, outputErr = s.actionService.Dispatch(ctx, toolCall.Function.Name, []byte(toolCall.Function.Arguments), metadata)
		}

		if outputErr != nil {
			output.Success = false
			output.Reason = outputErr.Error()
			logger.Warn("caused by callback", "name", toolCall.Function.Name, "err", outputErr)
		} else {
			output.Success = true
			output.Reason = ""
		}

		outputJson, err := json.Marshal(output)
		if err != nil {
			logger.Warn("failed to marshal output", "error", err)
		}

		toolOutput := openai.BetaThreadRunSubmitToolOutputsParamsToolOutput{
			ToolCallID: openai.F(toolCall.ID),
			Output:     openai.F(string(outputJson)),
		}

		outputs = append(outputs, toolOutput)
	}

	toolCallRun, err := s.openai.Beta.Threads.Runs.SubmitToolOutputs(ctx, run.ThreadID, run.ID, openai.BetaThreadRunSubmitToolOutputsParams{
		ToolOutputs: openai.F(outputs),
	})
	if err != nil {
		logger.Warn("failed to submit tool outputs", "error", err)
	}

	return toolCallRun, nil
}

func (s *server) doneAgent(
	ctx context.Context,
	thread *domain.Thread,
	agentWork *domain.AgentWork,
	actionWork *domain.ActionWork,
) (any, error) {
	tx := helpers.GetTx(ctx)

	if err := tx.Transaction(func(tx *gorm.DB) error {
		actionWork.Done = true
		if err := actionWork.Save(tx); err != nil {
			return err
		}

		agentWork.Status = domain.AgentStatusIdle
		if err := agentWork.Save(tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	logger.Info("done agent", "name", agentWork.Agent.Name, "action", actionWork.Action.Subject)

	var allActionWorks []domain.ActionWork
	if err := tx.
		InnerJoins("Action").
		InnerJoins("Action.Step", tx.Where(&domain.Step{SeqNo: thread.CurrentStepSeqNo})).
		Find(&allActionWorks, "thread_id = ?", thread.ID).Error; err != nil {
		return nil, err
	}

	if allDone := gog.Reduce(allActionWorks, true, func(acc bool, actionWork domain.ActionWork) bool {
		return acc && actionWork.Done
	}); !allDone {
		logger.Debug("not done step", "seq_no", thread.CurrentStepSeqNo)
		return nil, nil
	}

	nextStepSeqNo := thread.CurrentStepSeqNo + 1
	if len(thread.Mission.Steps) == 0 {
		return nil, errors.Errorf("no steps found")
	}
	isDone := nextStepSeqNo == thread.Mission.Steps[len(thread.Mission.Steps)-1].SeqNo+1
	hasNextStep := false
	for _, step := range thread.Mission.Steps {
		if step.SeqNo == nextStepSeqNo {
			hasNextStep = true
			break
		}
	}
	if !isDone && !hasNextStep {
		return nil, errors.Errorf("failed to find next step with seq no %d", nextStepSeqNo)
	}

	thread.AllDone = isDone
	if !isDone {
		thread.CurrentStepSeqNo = nextStepSeqNo
	}
	if err := thread.Save(tx); err != nil {
		return nil, err
	}

	logger.Info("done step", "seq_no", thread.CurrentStepSeqNo, "all_done", thread.AllDone)

	if hasNextStep {
		logger.Info("go to next step", "thread_id", thread.ID, "seq_no", nextStepSeqNo)
		helpers.On(ctx, helpers.EventTypeCompletedAction, func(ctx context.Context) {
			if err := s.run(ctx, thread.ID); err != nil {
				logger.Error("failed to run thread", "thread_id", thread.ID, "error", err)
			}
		})
	}

	return nil, nil
}
