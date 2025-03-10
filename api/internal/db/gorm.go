package db

import (
	"context"
	"github.com/habiliai/alice/api/config"
	"github.com/habiliai/alice/api/internal/di"
	"github.com/habiliai/alice/api/internal/mylog"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var (
	Key = di.NewKey()
)

func OpenDB(databaseUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	if db == nil {
		return errors.Errorf("db is nil")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return errors.Wrapf(err, "failed to get db")
	}
	if err := sqlDB.Close(); err != nil {
		return errors.Wrapf(err, "failed to close db")
	}

	return nil
}

func init() {
	di.Register(Key, func(c context.Context, env di.Env) (any, error) {
		logger, err := di.Get[*mylog.Logger](c, mylog.Key)
		if err != nil {
			return nil, err
		}

		cfg, err := di.Get[config.AliceConfig](c, config.AliceConfigKey)
		if err != nil {
			return nil, err
		}

		logger.Info("initialize database")
		db, err := OpenDB(cfg.DatabaseUrl)
		if err != nil {
			return nil, err
		}

		if env == di.EnvTest {
			if err := DropAll(db); err != nil {
				return nil, errors.Wrapf(err, "failed to drop database")
			}
			time.Sleep(500 * time.Millisecond)
		}
		if env == di.EnvTest || cfg.DatabaseAutoMigrate {
			if err := AutoMigrate(db); err != nil {
				return nil, errors.Wrapf(err, "failed to migrate database")
			}
		}

		go func() {
			<-c.Done()
			if err := CloseDB(db); err != nil {
				logger.Warn("failed to close database", "err", err)
			}
			logger.Info("database closed")
		}()

		return db, nil
	})
}
