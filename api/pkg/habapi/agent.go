package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) GetAgents(ctx context.Context, _ *emptypb.Empty) (*GetAgentsResponse, error) {
	tx := helpers.GetTx(ctx)

	var agents []domain.Agent
	if err := tx.Find(&agents).Error; err != nil {
		return nil, err
	}

	resp := &GetAgentsResponse{
		Agents:   make([]*Agent, 0, len(agents)),
		NumTotal: int32(len(agents)),
	}
	for _, agent := range agents {
		resp.Agents = append(resp.Agents, newAgentPbFromDb(&agent))
	}

	return resp, nil
}
