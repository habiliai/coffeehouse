package habapi

import (
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/mokiat/gog"
)

func newMissionPbFromDb(mission *domain.Mission) *Mission {
	m := &Mission{
		Id:   int32(mission.ID),
		Name: mission.Name,
		Steps: gog.Map(mission.Steps, func(task domain.Step) *Step {
			return newStepPbFromDb(&task)
		}),
	}

	memo := map[uint]struct{}{}
	for _, step := range mission.Steps {
		for _, action := range step.Actions {
			if _, ok := memo[action.AgentID]; !ok {
				memo[action.AgentID] = struct{}{}
				m.AgentIds = append(m.AgentIds, int32(action.AgentID))
			}
		}
	}

	return m
}

func newStepPbFromDb(step *domain.Step) *Step {
	return &Step{
		Id:    int32(step.ID),
		SeqNo: int32(step.SeqNo),
		Actions: gog.Map(step.Actions, func(action domain.Action) *Action {
			return newActionPbFromDb(&action)
		}),
	}
}

func newAgentPbFromDb(agent *domain.Agent) *Agent {
	return &Agent{
		Id:      int32(agent.ID),
		Name:    agent.Name,
		IconUrl: agent.IconUrl,
		Role:    agent.Role,
	}
}

func newActionPbFromDb(action *domain.Action) *Action {
	return &Action{
		Id:      int32(action.ID),
		Subject: action.Subject,
		Agent:   newAgentPbFromDb(&action.Agent),
	}
}
