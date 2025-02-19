package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"slices"
)

func (s *server) CreateThread(ctx context.Context, req *CreateThreadRequest) (*ThreadId, error) {
	tx := helpers.GetTx(ctx)

	var threadId uint
	if err := tx.Transaction(func(tx *gorm.DB) (err error) {
		var mission domain.Mission
		if err := tx.
			Preload("Steps").
			Preload("Steps.Actions").
			Preload("Steps.Actions.Agent").
			First(&mission, req.MissionId).Error; err != nil {
			return errors.Wrapf(err, "failed to find mission with id %d", req.MissionId)
		}

		var openaiThreadId string
		if thread, err := s.openai.Beta.Threads.New(ctx, openai.BetaThreadNewParams{}); err != nil {
			return errors.Wrap(err, "failed to create thread")
		} else {
			openaiThreadId = thread.ID
		}

		thread := domain.Thread{
			OpenaiThreadId:   openaiThreadId,
			MissionID:        mission.ID,
			Mission:          mission,
			CurrentStepSeqNo: 1,
			Data:             datatypes.NewJSONType(map[string]interface{}{}),
		}
		if err := thread.Save(tx); err != nil {
			return errors.Wrap(err, "failed to save thread")
		}

		slices.SortFunc(mission.Steps, func(a, b domain.Step) int {
			return a.SeqNo - b.SeqNo
		})

		var agents []*domain.Agent
		for _, step := range mission.Steps {
			for _, action := range step.Actions {
				actionWork := domain.ActionWork{
					ActionID: action.ID,
					ThreadID: thread.ID,

					Done: false,
				}
				if err := actionWork.Save(tx); err != nil {
					return errors.Wrap(err, "failed to save action work")
				}

				isDup := false
				for _, agent := range agents {
					if agent.ID == action.AgentID {
						isDup = true
						break
					}
				}
				if !isDup {
					agents = append(agents, &action.Agent)
				}
			}
		}

		for _, agent := range agents {
			agentWork := domain.AgentWork{
				AgentID:  agent.ID,
				ThreadID: thread.ID,
			}
			agentWork.Status = domain.AgentStatusIdle
			if err := agentWork.Save(tx); err != nil {
				return errors.Wrap(err, "failed to save agent work")
			}
		}

		threadId = thread.ID
		return nil
	}); err != nil {
		return nil, err
	}

	if threadId == 0 {
		return nil, errors.New("invalid thread id")
	}

	helpers.On(ctx, helpers.EventTypeEndTransaction, func(ctx context.Context) {
		logger.Debug("request to run thread", "threadId", threadId)
		if err := s.run(ctx, threadId); err != nil {
			logger.Error("failed to run thread", "threadId", threadId, "error", err)
		}
	})

	return &ThreadId{
		Id: int32(threadId),
	}, nil
}

func (s *server) GetThread(ctx context.Context, req *ThreadId) (*Thread, error) {
	openaiThread, thread, err := s.getThread(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res := &Thread{
		Id:               int32(thread.ID),
		MissionId:        int32(thread.MissionID),
		CurrentStepSeqNo: int32(thread.CurrentStepSeqNo),
		AllDone:          thread.AllDone,
	}

	res.Messages, err = s.getAllMessages(ctx, openaiThread.ID)
	if err != nil {
		return nil, err
	}

	if !thread.AllDone {
		return res, nil
	}

	if thread.AllDone {
		res.Result, err = s.generateResult(thread)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (s *server) GetThreadStatus(ctx context.Context, req *GetThreadStatusRequest) (*ThreadStatus, error) {
	_, thread, err := s.getThread(ctx, req.ThreadId)
	if err != nil {
		return nil, err
	}

	lastMessage, err := s.getLastMessage(ctx, thread.OpenaiThreadId)
	if err != nil {
		return nil, err
	}
	logger.Debug("last message", "lastMessage", lastMessage)

	lastMessageId := ""
	if lastMessage != nil {
		lastMessageId = lastMessage.Id
	}

	res := &ThreadStatus{}
	if req.LastMessageId != lastMessageId {
		res.HasNewMessage = true
	}

	return res, nil
}

func (s *server) getThread(ctx context.Context, id int32) (*openai.Thread, *domain.Thread, error) {
	tx := helpers.GetTx(ctx)

	var thread domain.Thread
	if err := tx.Preload("Mission").
		Preload("Mission.Steps").
		Preload("Mission.Steps.Actions").
		Preload("Mission.Steps.Actions.Agent").
		First(&thread, id).Error; err != nil {
		return nil, nil, errors.Wrapf(err, "failed to find thread with id %d", id)
	}

	openaiThread, err := s.openai.Beta.Threads.Get(ctx, thread.OpenaiThreadId)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get thread")
	}

	return openaiThread, &thread, nil
}
