package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
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
