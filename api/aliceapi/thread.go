package aliceapi

import (
	"context"
	_ "embed"
	"github.com/Masterminds/sprig/v3"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/agentruntime/thread"
	"github.com/habiliai/alice/api/domain"
	"github.com/habiliai/alice/api/internal/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"slices"
	"strings"
	"text/template"
	"time"
)

type AdditionalInstructionValues struct {
	Mission   string
	TodayDate string
	TodayDay  string
}

var (
	//go:embed data/instructions/thread.additional_instruction.md.tmpl
	additionalInstructions string

	additionalInstructionsTemplate = template.Must(template.New("thread_additional_instructions").Funcs(sprig.FuncMap()).Parse(additionalInstructions))
)

func newAdditionalInstructionValues(mission string) AdditionalInstructionValues {
	return AdditionalInstructionValues{
		Mission:   mission,
		TodayDate: time.Now().Format("2006-01-02"),
		TodayDay:  time.Now().Format("Monday"),
	}
}

func (s *server) CreateThread(ctx context.Context, req *CreateThreadRequest) (*ThreadId, error) {
	ctx, tx := db.OpenSession(ctx, s.db)

	var threadId uint
	var agentIds []uint32
	if err := tx.Transaction(func(tx *gorm.DB) (err error) {
		var mission domain.Mission
		if err := tx.
			Preload("Steps").
			Preload("Steps.Actions").
			First(&mission, req.MissionId).Error; err != nil {
			return errors.Wrapf(err, "failed to find mission with id %d", req.MissionId)
		}

		var instructionBuffer strings.Builder
		instructionValues := newAdditionalInstructionValues(mission.Name)
		if err := additionalInstructionsTemplate.Execute(&instructionBuffer, instructionValues); err != nil {
			return errors.Wrap(err, "failed to execute additional instructions template")
		}

		thr, err := s.threadManager.CreateThread(ctx, &thread.CreateThreadRequest{
			Instruction: instructionBuffer.String(),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create thread")
		}

		thd := domain.Thread{
			AgentRuntimeThreadId: thr.ThreadId,
			MissionID:            mission.ID,
			Mission:              mission,
			CurrentStepSeqNo:     1,
		}
		if err := thd.Save(tx); err != nil {
			return errors.Wrap(err, "failed to save thread")
		}

		slices.SortFunc(mission.Steps, func(a, b domain.Step) int {
			return a.SeqNo - b.SeqNo
		})

		for _, step := range mission.Steps {
			for _, action := range step.Actions {
				actionWork := domain.ActionWork{
					ActionID: action.ID,
					ThreadID: thd.ID,

					Done: false,
				}
				if err := actionWork.Save(tx); err != nil {
					return errors.Wrap(err, "failed to save action work")
				}

				isDup := false
				for _, agentId := range agentIds {
					if agentId == action.AgentId {
						isDup = true
						break
					}
				}
				if !isDup {
					agentIds = append(agentIds, action.AgentId)
				}
			}
		}

		for _, agentId := range agentIds {
			agentWork := domain.AgentWork{
				AgentId:  agentId,
				ThreadID: thd.ID,
			}
			agentWork.Status = domain.AgentStatusIdle
			if err := agentWork.Save(tx); err != nil {
				return errors.Wrap(err, "failed to save agent work")
			}
		}

		threadId = thd.ID

		return nil
	}); err != nil {
		return nil, err
	}

	if err := s.run(ctx, threadId, agentIds); err != nil {
		s.logger.Error("failed to run thread", "threadId", threadId, "error", err)
	}

	return &ThreadId{
		Id: int32(threadId),
	}, nil
}

func (s *server) GetThread(ctx context.Context, req *ThreadId) (*Thread, error) {
	ctx, tx := db.OpenSession(ctx, s.db)
	thd, err := s.getThread(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var actions []domain.Action
	if err := tx.
		Model(&domain.Action{}).
		InnerJoins("Step", &domain.Step{MissionID: thd.MissionID}).
		Find(&actions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find action works")
	}

	agent2Actions := map[string]*struct {
		Action *domain.Action
		Agent  *agent.Agent
		Done   bool
	}{}

	for _, a := range actions {
		ag, err := s.agentManager.GetAgent(ctx, &agent.GetAgentRequest{
			AgentId: a.AgentId,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		agent2Actions[ag.Name] = &struct {
			Action *domain.Action
			Agent  *agent.Agent
			Done   bool
		}{
			Action: &a,
			Agent:  ag,
			Done:   false,
		}
	}

	messages, err := s.getAllMessages(ctx, thd.AgentRuntimeThreadId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find messages")
	}

	for _, msg := range messages {
		a, ok := agent2Actions[msg.Sender]
		if !ok {
			continue
		}
		for _, toolCall := range msg.ToolCalls {
			if toolCall.Name == "done_agent" {
				a.Done = true
			}
		}
	}

	res := &Thread{
		Id:               int32(thd.ID),
		MissionId:        int32(thd.MissionID),
		CurrentStepSeqNo: int32(thd.CurrentStepSeqNo),
		AllDone:          thd.AllDone,
	}

	for _, a := range agent2Actions {
		action, err := s.newActionPb(ctx, a.Action, a.Agent)
		if err != nil {
			return nil, err
		}
		res.ActionWorks = append(res.ActionWorks, &ActionWork{
			Action: action,
			Done:   a.Done,
		})
	}

	res.Messages, err = s.convertAllMessagesFromAgentRuntime(ctx, messages)
	if err != nil {
		return nil, err
	}

	if thd.AllDone {
		res.Result, err = s.generateResult(ctx, thd)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (s *server) GetThreadStatus(ctx context.Context, req *GetThreadStatusRequest) (*ThreadStatus, error) {
	thread, err := s.getThread(ctx, req.ThreadId)
	if err != nil {
		return nil, err
	}

	lastMessage, err := s.getLastMessage(ctx, thread.AgentRuntimeThreadId)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("last message", "lastMessage", lastMessage)

	lastMessageId := uint32(0)
	if lastMessage != nil {
		lastMessageId = lastMessage.Id
	}

	res := &ThreadStatus{}
	if req.LastMessageId != lastMessageId {
		res.HasNewMessage = true
	}

	return res, nil
}

func (s *server) getThread(ctx context.Context, id int32) (*domain.Thread, error) {
	ctx, tx := db.OpenSession(ctx, s.db)

	var thread domain.Thread
	if err := tx.Preload("Mission").
		Preload("Mission.Steps").
		Preload("Mission.Steps.Actions").
		First(&thread, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find thread with id %d", id)
	}

	return &thread, nil
}
