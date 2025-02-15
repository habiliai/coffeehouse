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
				Name:        "scheduler",
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

			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&flags.reset, "reset", false, "Reset the database")

	return cmd
}
