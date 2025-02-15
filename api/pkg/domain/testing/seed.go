package domaintest

import (
	"github.com/habiliai/habiliai/api/pkg/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Seed struct {
	Agents   []*domain.Agent
	Missions []*domain.Mission
	Tasks    []*domain.Task
}

func SeedForTest(tx *gorm.DB) (r Seed, err error) {
	mission := &domain.Mission{
		Name: "Mission 1",
	}
	if err := mission.Save(tx); err != nil {
		logger.Warn("failed to save mission")
	}
	r.Missions = append(r.Missions, mission)

	r.Tasks = []*domain.Task{
		{
			SeqNo:     1,
			MissionID: mission.ID,
		},
		{
			SeqNo:     2,
			MissionID: mission.ID,
		},
		{
			SeqNo:     3,
			MissionID: mission.ID,
		},
		{
			SeqNo:     4,
			MissionID: mission.ID,
		},
		{
			SeqNo:     5,
			MissionID: mission.ID,
		},
	}
	for _, task := range r.Tasks {
		if err := task.Save(tx); err != nil {
			logger.Warn("failed to save task")
		}
	}

	r.Agents = []*domain.Agent{
		{
			Name:        "engineer",
			AssistantId: "asst_1Ov3IAylZU7Z9ansD7ZBWufs",
			IconUrl:     "https://img.logo.dev/github.com",
			Missions: []Mission{
				*r.Missions[0],
			},
		},
		{
			Name:        "designer",
			AssistantId: "asst_1Ov3IAylZU7Z9ansD7ZBWufs",
			IconUrl:     "https://img.logo.dev/facebook.com",
			Missions: []Mission{
				*r.Missions[0],
			},
		},
		{
			Name:        "manager",
			AssistantId: "asst_1Ov3IAylZU7Z9ansD7ZBWufs",
			IconUrl:     "https://img.logo.dev/netflix.com",
			Missions: []Mission{
				*r.Missions[0],
			},
		},
	}

	for _, agent := range r.Agents {
		if err := agent.Save(tx); err != nil {
			logger.Warn("failed to save agent")
		}
	}

	return
}
