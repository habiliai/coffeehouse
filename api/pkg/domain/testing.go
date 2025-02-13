package domain

import "gorm.io/gorm"

func SeedForTest(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		mission := Mission{
			Name: "Mission 1",
		}
		if err := mission.Save(tx); err != nil {
			return err
		}

		tasks := []Task{
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
		for _, task := range tasks {
			if err := task.Save(tx); err != nil {
				return err
			}
		}

		agents := []Agent{
			{
				Name:        "engineer",
				AssistantId: "asst_JwxSg8eflu0qGI0ZQ3tKSaCi",
				IconUrl:     "https://img.logo.dev/github.com",
				Missions: []Mission{
					mission,
				},
			},
			{
				Name:        "designer",
				AssistantId: "asst_JwxSg8eflu0qGI0ZQ3tKSaCi",
				IconUrl:     "https://img.logo.dev/facebook.com",
				Missions: []Mission{
					mission,
				},
			},
			{
				Name:        "manager",
				AssistantId: "asst_JwxSg8eflu0qGI0ZQ3tKSaCi",
				IconUrl:     "https://img.logo.dev/netflix.com",
				Missions: []Mission{
					mission,
				},
			},
		}

		for _, agent := range agents {
			if err := agent.Save(tx); err != nil {
				return err
			}
		}
		return nil
	})
}
