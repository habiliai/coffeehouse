package services

import (
	"github.com/habiliai/alice/api/pkg/constants"
	"github.com/habiliai/alice/api/pkg/digo"
	"github.com/habiliai/alice/api/pkg/domain"
	aflog "github.com/habiliai/alice/api/pkg/log"
	"github.com/habiliai/alice/api/pkg/util"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const ServiceKeyDB = "db"

func init() {
	digo.Register(ServiceKeyDB, func(ctx *digo.Container) (any, error) {
		switch ctx.Env {
		case digo.EnvProd:
			db, err := util.NewDBFromConfig(ctx, ctx.Config.DB)
			logger.Debug("new", "db", db)
			if err != nil {
				return nil, err
			}
			go func() {
				<-ctx.Done()
				if err := util.CloseDB(db); err != nil {
					logger.Warn("failed to close database", aflog.Err(err))
				}
				logger.Info("database closed")
			}()

			return db, nil
		case digo.EnvTest:
			db, err := gorm.Open(postgres.Open("postgres://habiliai:habiliai@localhost:5432/test?search_path="+constants.SchemaName), &gorm.Config{})
			if err != nil {
				return nil, err
			}
			go func() {
				<-ctx.Done()
				if err := util.CloseDB(db); err != nil {
					logger.Warn("failed to close database", aflog.Err(err))
				}
				logger.Info("database closed")
			}()

			if err := domain.DropAll(db); err != nil {
				logger.Warn("failed to drop all tables", aflog.Err(err))
			}
			time.Sleep(500 * time.Millisecond)
			if err := domain.AutoMigrate(db); err != nil {
				return nil, err
			}
			return db, nil
		default:
			return nil, errors.New("unknown env")
		}
	})
}
