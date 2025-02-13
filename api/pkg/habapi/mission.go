package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/pkg/errors"
)

func (s *server) GetMissions(ctx context.Context, req *GetMissionsRequest) (*GetMissionsResponse, error) {
	tx := helpers.GetTx(ctx)

	stmt := tx.Model(&domain.Mission{})
	var numTotal int64
	if err := stmt.Count(&numTotal).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count missions")
	}

	var missions []domain.Mission
	if err := stmt.Order("id ASC").Find(&missions).Error; err != nil {
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
