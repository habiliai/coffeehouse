package habapi

import "github.com/habiliai/habiliai/api/pkg/domain"

func newMissionPbFromDb(mission *domain.Mission) *Mission {
	return &Mission{
		Id:   int32(mission.ID),
		Name: mission.Name,
	}
}

func newAgentPbFromDb(agent *domain.Agent) *Agent {
	return &Agent{
		Id:      int32(agent.ID),
		Name:    agent.Name,
		IconUrl: agent.IconUrl,
	}
}
