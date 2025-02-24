package util

import (
	"context"
	"github.com/habiliai/alice/api/pkg/config"
	"github.com/habiliai/alice/api/pkg/domain"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func NewDBFromConfig(ctx context.Context, config config.DBConfig) (*gorm.DB, error) {
	return NewDBFromConfigAndNamingStrategy(ctx, config, schema.NamingStrategy{})
}

func NewDBFromConfigAndNamingStrategy(ctx context.Context, config config.DBConfig, strategy schema.NamingStrategy) (db *gorm.DB, err error) {
	if strategy.IdentifierMaxLength == 0 {
		strategy.IdentifierMaxLength = 64
	}

	var pingTimeout time.Duration
	if config.PingTimeout != "" {
		pingTimeout, err = time.ParseDuration(config.PingTimeout)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse ping timeout")
		}
	}

	uri := config.GetURI()
	logger.Debug("connecting database", "uri", uri)
	db, err = gorm.Open(postgres.Open(uri), &gorm.Config{
		TranslateError: true,
		NamingStrategy: strategy,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open db")
	}
	logger.Debug("opened db", "db", db, "strategy", strategy)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get db")
	}

	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, errors.Wrapf(err, "failed to ping db")
	}
	logger.Debug("db ping is ok")

	if config.AutoMigration {
		if err := domain.AutoMigrate(db); err != nil {
			return nil, errors.Wrapf(err, "failed to auto migrate db")
		}
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	if value, err := time.ParseDuration(config.ConnMaxLifetime); err != nil {
		return nil, errors.Wrapf(err, "failed to parse conn max lifetime")
	} else {
		sqlDB.SetConnMaxLifetime(value)
	}
	logger.Debug("succeeded to create gorm db")

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return errors.Wrapf(err, "failed to get db")
	}
	if err := sqlDB.Close(); err != nil {
		return errors.Wrapf(err, "failed to close db")
	}

	logger.Info("closing db")
	return nil
}
