package habapi

import (
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/mokiat/gog"
)

func newMissionPbFromDb(mission *domain.Mission) *Mission {
	return &Mission{
		Id:   int32(mission.ID),
		Name: mission.Name,
		Agents: gog.Map(mission.Agents, func(agent domain.Agent) *Agent {
			return newAgentPbFromDb(&agent)
		}),
	}
}

func newAgentPbFromDb(agent *domain.Agent) *Agent {
	return &Agent{
		Id:      int32(agent.ID),
		Name:    agent.Name,
		IconUrl: agent.IconUrl,
	}
}
