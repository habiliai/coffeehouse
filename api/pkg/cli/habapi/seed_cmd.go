package habapi

import (
	"github.com/habiliai/habiliai/api/pkg/digo"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/services"
	"github.com/pkg/errors"
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
				IncludeQuestionIntent: true,
			}
			moderator.Save(db)
			suggester := domain.Agent{
				Name:        "suggester",
				AssistantId: "asst_YHXPoVMv8oD3kzBO2xO8jIxL",
				IconUrl:     "https://img.logo.dev/facebook.com",
			}
			suggester.Save(db)
			scheduler := domain.Agent{
				Name:        "scheduler1",
				AssistantId: "asst_NesAuSvy09nlv7gI3c2Bcx61",
				IconUrl:     "https://img.logo.dev/google.com",
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
				return errors.Wrapf(err, "failed to create mission")
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
			}
			communityManager.Save(db)
			locationRecommender := domain.Agent{
				Name:                  "loc",
				IconUrl:               "https://img.logo.dev/facebook.com",
				AssistantId:           "asst_7B77ZJoBXpia1QB0E5G8yOnG",
				IncludeQuestionIntent: true,
			}
			locationRecommender.Save(db)
			commentSpecialist := domain.Agent{
				Name:        "commenter",
				IconUrl:     "https://img.logo.dev/google.com",
				AssistantId: "asst_Duu1WXwmYwPokx3n3JYIVvRo",
			}
			commentSpecialist.Save(db)
			scheduleManager := domain.Agent{
				Name:        "scheduler2",
				IconUrl:     "https://img.logo.dev/whatsapp.com",
				AssistantId: "asst_jLfJ39wLvYRKiy21Iirq6fIl",
			}
			scheduleManager.Save(db)
			snsManager := domain.Agent{
				Name:        "sns",
				IconUrl:     "https://img.logo.dev/instagram.com",
				AssistantId: "asst_OjOnv01dbmaGMa200m9bkDpm",
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
				return errors.Wrapf(err, "failed to create mission")
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&flags.reset, "reset", false, "Reset the database")

	return cmd
}
