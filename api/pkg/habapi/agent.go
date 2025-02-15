package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/mokiat/gog"
	"github.com/pkg/errors"
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

func (s *server) GetAgentsStatus(ctx context.Context, req *ThreadId) (*AgentsStatus, error) {
	_, thread, err := s.getThread(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var works []domain.AgentWork
	if err := helpers.GetTx(ctx).
		Preload("Agent").
		Order("id ASC").
		Find(&works, "thread_id = ?", thread.ID).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find agent works")
	}

	return &AgentsStatus{
		Works: gog.Map(works, func(w domain.AgentWork) *AgentWork {
			res := &AgentWork{
				Agent: newAgentPbFromDb(&w.Agent),
			}

			switch w.Status {
			case domain.AgentStatusWorking:
				res.Status = AgentWork_WORKING
			case domain.AgentStatusIdle:
				res.Status = AgentWork_IDLE
			case domain.AgentStatusWaiting:
				res.Status = AgentWork_WAITING
			}

			return res
		}),
	}, nil
}
