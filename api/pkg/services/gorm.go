package services

import (
	"github.com/habiliai/habiliai/api/pkg/digo"
	"github.com/habiliai/habiliai/api/pkg/domain"
	aflog "github.com/habiliai/habiliai/api/pkg/log"
	"github.com/habiliai/habiliai/api/pkg/util"
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
				logger.Info("closing database")
				if err := util.CloseDB(db); err != nil {
					logger.Warn("failed to close database", aflog.Err(err))
				}
			}()

			return db, nil
		case digo.EnvTest:
			db, err := gorm.Open(postgres.Open("postgres://afuser:afuser@localhost:6543/test?search_path=agentfather"), &gorm.Config{})
			if err != nil {
				return nil, err
			}
			go func() {
				<-ctx.Done()
				logger.Info("closing database")
				if err := util.CloseDB(db); err != nil {
					logger.Warn("failed to close database", aflog.Err(err))
				}
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
