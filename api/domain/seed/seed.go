package seed

import (
	"context"
	_ "embed"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/alice/api/domain"
	"github.com/habiliai/alice/api/internal/agentruntimeclient"
	"github.com/habiliai/alice/api/internal/db"
	"github.com/habiliai/alice/api/internal/di"
	"github.com/habiliai/alice/api/internal/mylog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	//go:embed data/templates/mission1_result.md.tmpl
	mission1ResultTemplate string
)

func Seed(ctx context.Context, reset bool) error {
	dbInstance := di.MustGet[*gorm.DB](ctx, db.Key)
	logger := di.MustGet[*mylog.Logger](ctx, mylog.Key)

	if reset {
		if err := db.DropAll(dbInstance); err != nil {
			return err
		}
		if err := db.AutoMigrate(dbInstance); err != nil {
			return err
		}
	}

	agentManager := agentruntimeclient.GetAgentManager(ctx)
	sunny, err := agentManager.GetAgentByName(ctx, &agent.GetAgentByNameRequest{
		Name: "sunny",
	})
	if err != nil {
		return errors.WithStack(err)
	}
	eric, err := agentManager.GetAgentByName(ctx, &agent.GetAgentByNameRequest{
		Name: "eric",
	})
	if err != nil {
		return errors.WithStack(err)
	}

	mission := domain.Mission{
		Name:           "Organize a community meal gathering near Hong Kong CEC this weekend.",
		ResultTemplate: mission1ResultTemplate,
		Steps: []domain.Step{
			{
				SeqNo: 1,
				Actions: []domain.Action{
					{
						AgentId: sunny.Id,
						Subject: "Check the weather forecast for this weekend",
					},
					{
						AgentId: eric.Id,
						Subject: "Find a suitable location for the gathering",
					},
				},
			},
			//{
			//	SeqNo: 2,
			//	Actions: []Action{
			//		{
			//			AgentId: julia.ID,
			//			Subject: "Recommend a meeting place",
			//		},
			//	},
			//},
			//{
			//	SeqNo: 3,
			//	Actions: []Action{
			//		{
			//			AgentId: amelia.ID,
			//			Subject: "Schedule the gathering",
			//		},
			//	},
			//},
			//{
			//	SeqNo: 4,
			//	Actions: []Action{
			//		{
			//			AgentId: stella.ID,
			//			Subject: "Write a comment to invite people",
			//		},
			//		{
			//			AgentId: nolan.ID,
			//			Subject: "Post the event on X",
			//		},
			//	},
			//},
		},
	}

	if err := dbInstance.Create(&mission).Error; err != nil {
		logger.Warn("failed to create mission", "err", err)
	}

	return nil
}
