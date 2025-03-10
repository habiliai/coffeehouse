package aliceapi

import (
	"context"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/alice/api/domain"
	"github.com/habiliai/alice/api/internal/db"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"slices"
)

func (s *server) GetAgents(ctx context.Context, _ *emptypb.Empty) (*GetAgentsResponse, error) {
	ctx, tx := db.OpenSession(ctx, s.db)

	var actions []domain.Action
	if err := tx.
		Find(&actions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find agent ids")
	}

	agentIds := make([]uint32, 0, len(actions))
	for _, action := range actions {
		if action.AgentId == 0 {
			continue
		}
		agentIds = append(agentIds, action.AgentId)
	}

	// Remove duplicates
	agentIds = slices.Compact(agentIds)

	agents := make([]*agent.Agent, 0, len(agentIds))
	for _, agentId := range agentIds {
		ag, err := s.agentManager.GetAgent(ctx, &agent.GetAgentRequest{
			AgentId: agentId,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find agent with id %d", agentId)
		}
		agents = append(agents, ag)
	}

	resp := &GetAgentsResponse{
		Agents:   make([]*Agent, 0, len(agents)),
		NumTotal: int32(len(agents)),
	}
	for _, agent := range agents {
		resp.Agents = append(resp.Agents, newAgentPbFromAgentRuntime(agent))
	}

	return resp, nil
}

func (s *server) GetAgentsStatus(ctx context.Context, req *ThreadId) (*AgentsStatus, error) {
	ctx, tx := db.OpenSession(ctx, s.db)

	thread, err := s.getThread(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var works []domain.AgentWork
	if err := tx.
		Order("agent_id ASC").
		Find(&works, "thread_id = ?", thread.ID).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find agent works")
	}

	res := AgentsStatus{
		Works: make([]*AgentWork, 0, len(works)),
	}
	for _, w := range works {
		ag, err := s.agentManager.GetAgent(ctx, &agent.GetAgentRequest{
			AgentId: w.AgentId,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find agent with id %d", w.AgentId)
		}
		agentWork := &AgentWork{
			Agent: newAgentPbFromAgentRuntime(ag),
		}

		switch w.Status {
		case domain.AgentStatusWorking:
			agentWork.Status = AgentWork_WORKING
		case domain.AgentStatusIdle:
			agentWork.Status = AgentWork_IDLE
		case domain.AgentStatusWaiting:
			agentWork.Status = AgentWork_WAITING
		}

		res.Works = append(res.Works, agentWork)
	}

	return &res, nil
}
