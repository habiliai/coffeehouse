package aliceapi

import (
	"context"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/alice/api/domain"
	"github.com/pkg/errors"
)

func (s *server) newMissionPbFromDb(ctx context.Context, mission *domain.Mission) (*Mission, error) {
	m := &Mission{
		Id:    int32(mission.ID),
		Name:  mission.Name,
		Steps: make([]*Step, 0, len(mission.Steps)),
	}

	for _, step := range mission.Steps {
		st, err := s.newStepPbFromDb(ctx, &step)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create step pb")
		}
		m.Steps = append(m.Steps, st)
	}

	memo := map[uint32]struct{}{}
	for _, step := range mission.Steps {
		for _, action := range step.Actions {
			if _, ok := memo[action.AgentId]; !ok {
				memo[action.AgentId] = struct{}{}
				m.AgentIds = append(m.AgentIds, int32(action.AgentId))
			}
		}
	}

	return m, nil
}

func (s *server) newStepPbFromDb(ctx context.Context, step *domain.Step) (*Step, error) {
	st := &Step{
		Id:      int32(step.ID),
		SeqNo:   int32(step.SeqNo),
		Actions: make([]*Action, 0, len(step.Actions)),
	}

	for _, action := range step.Actions {
		actionPb, err := s.newActionPb(ctx, &action, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create action pb")
		}
		st.Actions = append(st.Actions, actionPb)
	}

	return st, nil
}

func newAgentPbFromAgentRuntime(agent *agent.Agent) *Agent {
	iconUrl := agent.Metadata["iconUrl"]
	return &Agent{
		Id:      int32(agent.Id),
		Name:    agent.Name,
		IconUrl: iconUrl,
		Role:    agent.Role,
	}
}

func (s *server) newActionPb(ctx context.Context, action *domain.Action, ag *agent.Agent) (*Action, error) {
	var err error
	if ag == nil {
		ag, err = s.agentManager.GetAgent(ctx, &agent.GetAgentRequest{
			AgentId: action.AgentId,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find agent with id %d", action.AgentId)
		}
	}

	return &Action{
		Id:      int32(action.ID),
		Subject: action.Subject,
		Agent:   newAgentPbFromAgentRuntime(ag),
	}, nil
}
