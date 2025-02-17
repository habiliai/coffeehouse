package habapi

import (
	"github.com/habiliai/habiliai/api/pkg/digo"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/services"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func (c *cli) newSeedCmd() *cobra.Command {
	flags := &struct {
		reset bool
	}{}

	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed the database",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			if err := c.ReadInConfig(); err != nil {
				return err
			}

			container := digo.NewContainer(ctx, digo.EnvProd, &c.cfg)
			db, err := digo.Get[*gorm.DB](container, services.ServiceKeyDB)
			if err != nil {
				return err
			}

			if flags.reset {
				if err := domain.DropAll(db); err != nil {
					return err
				}
				if err := domain.AutoMigrate(db); err != nil {
					return err
				}
			}

			moderator := domain.Agent{
				Name:                  "moderator",
				AssistantId:           "asst_aggf9nmtEM77Qy3niFq64uBK",
				IconUrl:               "https://img.logo.dev/github.com",
				Role:                  "Identify the user’s needs and propose recommended keywords as options.",
				IncludeQuestionIntent: true,
			}
			moderator.Save(db)
			suggester := domain.Agent{
				Name:        "suggester",
				AssistantId: "asst_YHXPoVMv8oD3kzBO2xO8jIxL",
				IconUrl:     "https://img.logo.dev/facebook.com",
				Role:        "The AI serves as a conference manager that recommends relevant schedules and events based on user interests, providing easily selectable options for registration and applications.",
			}
			suggester.Save(db)
			scheduler := domain.Agent{
				Name:        "scheduler1",
				AssistantId: "asst_NesAuSvy09nlv7gI3c2Bcx61",
				IconUrl:     "https://img.logo.dev/google.com",
				Role:        "The AI acts as a schedule manager that identifies suitable events from the user’s calendar availability and, upon confirmation, automatically registers the chosen schedule.",
			}
			scheduler.Save(db)

			mission := domain.Mission{
				Name: "컨퍼런스 일정 중에 갈만한 것을 예약해줘",
				Steps: []domain.Step{
					{
						SeqNo: 1,
						Actions: []domain.Action{
							{
								Subject: "니즈 파악(추천 키워드 선택지 제안)",
								Agent:   moderator,
							},
						},
					},
					{
						SeqNo: 2,
						Actions: []domain.Action{
							{
								Subject: "일정, 이벤트 추천",
								Agent:   suggester,
							},
						},
					},
					{
						SeqNo: 3,
						Actions: []domain.Action{
							{
								Subject: "적절한 일정 찾아 신청",
								Agent:   scheduler,
							},
						},
					},
				},
			}

			if err := db.Create(&mission).Error; err != nil {
				logger.Warn("failed to create mission", "err", err)
			}

			weatherForecaster := domain.Agent{
				Name:        "weather",
				AssistantId: "asst_Rho49KGmpl1IkiVWtfuNDd4i",
				IconUrl:     "https://img.logo.dev/twitter.com",
			}
			weatherForecaster.Save(db)
			communityManager := domain.Agent{
				Name:                  "comm",
				IconUrl:               "https://img.logo.dev/github.com",
				IncludeQuestionIntent: true,
				AssistantId:           "asst_e08WipDCbvGTOVlFxiMjPa10",
				Role:                  "The AI acts as a weather forecaster, checking the upcoming weekend weather at the specified location (e.g., HKCEC) using OpenWeatherMap to recommend suitable dates.",
			}
			communityManager.Save(db)
			locationRecommender := domain.Agent{
				Name:                  "loc",
				IconUrl:               "https://img.logo.dev/facebook.com",
				AssistantId:           "asst_7B77ZJoBXpia1QB0E5G8yOnG",
				IncludeQuestionIntent: true,
				Role:                  "The AI acts as a community manager that identifies a user’s food preferences via LLM-based interactions and provides an explanatory summary of their tastes.",
			}
			locationRecommender.Save(db)
			commentSpecialist := domain.Agent{
				Name:        "commenter",
				IconUrl:     "https://img.logo.dev/google.com",
				AssistantId: "asst_Duu1WXwmYwPokx3n3JYIVvRo",
				Role:        "The AI serves as a location manager that retrieves open place listings from a predefined database and recommends suitable venues based on user preferences, including details such as addresses, tags, ratings, and corresponding map links.",
			}
			commentSpecialist.Save(db)
			scheduleManager := domain.Agent{
				Name:        "scheduler2",
				IconUrl:     "https://img.logo.dev/whatsapp.com",
				AssistantId: "asst_jLfJ39wLvYRKiy21Iirq6fIl",
				Role:        "The AI functions as a ‘message specialist’ that crafts romantically-toned yet professional inquiries for business contacts, seamlessly incorporating schedule, time, location, and weather details into the message.",
			}
			scheduleManager.Save(db)
			snsManager := domain.Agent{
				Name:        "sns",
				IconUrl:     "https://img.logo.dev/instagram.com",
				AssistantId: "asst_OjOnv01dbmaGMa200m9bkDpm",
				Role:        "The AI acts as a schedule manager that confirms and registers event details in the user’s calendar, then provides a prompt confirmation of the successful addition.",
			}
			snsManager.Save(db)

			mission = domain.Mission{
				Name: "Organize a community meal gathering near Hong Kong CEC this weekend.",
				Steps: []domain.Step{
					{
						SeqNo: 1,
						Actions: []domain.Action{
							{
								AgentID: weatherForecaster.ID,
								Subject: "Check the weather forecast for this weekend",
							},
							{
								AgentID: communityManager.ID,
								Subject: "Find a suitable location for the gathering",
							},
						},
					},
					{
						SeqNo: 2,
						Actions: []domain.Action{
							{
								AgentID: locationRecommender.ID,
								Subject: "Recommend a meeting place",
							},
						},
					},
					{
						SeqNo: 3,
						Actions: []domain.Action{
							{
								AgentID: commentSpecialist.ID,
								Subject: "Write a comment to invite people",
							},
						},
					},
					{
						SeqNo: 4,
						Actions: []domain.Action{
							{
								AgentID: scheduleManager.ID,
								Subject: "Schedule the gathering",
							},
							{
								AgentID: snsManager.ID,
								Subject: "Post the event on X",
							},
						},
					},
				},
			}

			if err := db.Create(&mission).Error; err != nil {
				logger.Warn("failed to create mission", "err", err)
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&flags.reset, "reset", false, "Reset the database")

	return cmd
}
