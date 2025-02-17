package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/mokiat/gog"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
)

func (s *server) GetMissions(ctx context.Context, _ *emptypb.Empty) (*GetMissionsResponse, error) {
	tx := helpers.GetTx(ctx)

	stmt := tx.Model(&domain.Mission{})
	var numTotal int64
	if err := stmt.Count(&numTotal).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count missions")
	}

	var missions []domain.Mission
	if err := stmt.Preload(clause.Associations).
		Preload("Steps.Actions").
		Preload("Steps.Actions.Agent").
		Order("id ASC").
		Find(&missions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find missions")
	}

	resp := &GetMissionsResponse{
		Missions: make([]*Mission, len(missions)),
		NumTotal: int32(numTotal),
	}

	for i := range missions {
		resp.Missions[i] = newMissionPbFromDb(&missions[i])
	}

	return resp, nil
}

func (s *server) GetMission(ctx context.Context, req *MissionId) (*Mission, error) {
	tx := helpers.GetTx(ctx)

	var mission domain.Mission
	err := tx.
		Preload(clause.Associations).
		Preload("Steps.Actions").
		Preload("Steps.Actions.Agent").
		First(&mission, req.Id).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find mission by id %d", req.Id)
	}

	return newMissionPbFromDb(&mission), nil
}

func (s *server) GetMissionStepStatus(ctx context.Context, req *GetMissionStepStatusRequest) (*MissionStepStatus, error) {
	tx := helpers.GetTx(ctx)
	_, thread, err := s.getThread(ctx, req.ThreadId)
	if err != nil {
		return nil, err
	}

	var works []domain.ActionWork
	if err := helpers.GetTx(ctx).
		Preload("Action").
		Preload("Action.Agent").
		Preload("Action.Step").
		Find(
			&works,
			"thread_id = ? AND Action__Step.seq_no = ? AND Action.mission_id = ?",
			thread.ID, thread.CurrentStepSeqNo, thread.MissionID,
		).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find action works")
	}

	var step domain.Step
	if err := tx.
		First(&step, "mission_id = ? AND seq_no = ?", thread.MissionID, thread.CurrentStepSeqNo).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find step")
	}

	return &MissionStepStatus{
		Step: newStepPbFromDb(&step),
		ActionWorks: gog.Map(works, func(w domain.ActionWork) *ActionWork {
			res := &ActionWork{
				Action: newActionPbFromDb(&w.Action),
				Done:   w.Done,
			}
			if w.Error != "" {
				res.Error = &w.Error
			}

			return res
		}),
	}, nil
}
