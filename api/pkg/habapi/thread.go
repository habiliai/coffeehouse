package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *server) CreateThread(ctx context.Context, req *CreateThreadRequest) (*ThreadId, error) {
	tx := helpers.GetTx(ctx)

	var thread *openai.Thread
	if err := tx.Transaction(func(tx *gorm.DB) (err error) {
		var mission domain.Mission
		if err := tx.First(&mission, req.MissionId).Error; err != nil {
			return errors.Wrapf(err, "failed to find mission with id %d", req.MissionId)
		}

		thread, err = s.openai.Beta.Threads.New(ctx, openai.BetaThreadNewParams{})
		if err != nil {
			return errors.Wrap(err, "failed to create thread")
		}

		threadData := NewThreadData(s.openai, thread.ID, thread.Metadata)
		threadData.SetMissionId(uint(mission.ID))
		if err := threadData.Save(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &ThreadId{
		Id: thread.ID,
	}, nil
}

func (s *server) DeleteThread(ctx context.Context, req *ThreadId) (*emptypb.Empty, error) {
	if _, err := s.openai.Beta.Threads.Delete(ctx, req.Id); err != nil {
		return nil, errors.Wrapf(err, "failed to delete thread with id %s", req.Id)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) GetThread(ctx context.Context, req *ThreadId) (*Thread, error) {
	thread, threadData, err := s.getThread(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res := &Thread{
		Id:        thread.ID,
		MissionId: int32(threadData.GetMissionId()),
	}

	res.Messages, err = s.getAllMessages(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) getThread(ctx context.Context, id string) (*openai.Thread, *ThreadData, error) {
	thread, err := s.openai.Beta.Threads.Get(ctx, id)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get thread")
	}

	threadData := NewThreadData(s.openai, id, thread.Metadata)

	return thread, threadData, nil
}
