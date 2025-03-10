package db

import (
	"github.com/habiliai/alice/api/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + SchemaName).Error; err != nil {
		return errors.Wrapf(err, "failed to create schema")
	}

	return errors.WithStack(db.AutoMigrate(
		&domain.Mission{},
		&domain.Thread{},
		&domain.Step{},
		&domain.Action{},
		&domain.ActionWork{},
		&domain.AgentWork{},
	))
}

func DropAll(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&domain.AgentWork{},
		&domain.ActionWork{},
		&domain.Action{},
		&domain.Step{},
		&domain.Thread{},
		&domain.Mission{},
	)
}
