package domain

import "gorm.io/gorm"

type Seed struct {
	Agents   []*Agent
	Missions []*Mission
	Tasks    []*Task
}

func SeedForTest(db *gorm.DB) (r Seed, err error) {
	if err = db.Transaction(func(tx *gorm.DB) error {
		mission := &Mission{
			Name: "Mission 1",
		}
		if err := mission.Save(tx); err != nil {
			return err
		}
		r.Missions = append(r.Missions, mission)

		r.Tasks = []*Task{
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
				return err
			}
		}

		r.Agents = []*Agent{
			{
				Name:        "engineer",
				AssistantId: "asst_JwxSg8eflu0qGI0ZQ3tKSaCi",
				IconUrl:     "https://img.logo.dev/github.com",
				Missions: []Mission{
					*r.Missions[0],
				},
			},
			{
				Name:        "designer",
				AssistantId: "asst_JwxSg8eflu0qGI0ZQ3tKSaCi",
				IconUrl:     "https://img.logo.dev/facebook.com",
				Missions: []Mission{
					*r.Missions[0],
				},
			},
			{
				Name:        "manager",
				AssistantId: "asst_JwxSg8eflu0qGI0ZQ3tKSaCi",
				IconUrl:     "https://img.logo.dev/netflix.com",
				Missions: []Mission{
					*r.Missions[0],
				},
			},
		}

		for _, agent := range r.Agents {
			if err := agent.Save(tx); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return
	}

	return
}
