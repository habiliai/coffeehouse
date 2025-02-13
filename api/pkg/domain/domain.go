package domain

import (
	"fmt"
	"github.com/habiliai/habiliai/api/pkg/constants"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, constants.SchemaName)).Error; err != nil {
		return errors.Wrapf(err, "failed to create schema")
	}

	if err := errors.Wrapf(db.
		AutoMigrate(
			&Agent{},
			&Mission{},
			&Task{},
		), "failed to auto migrate"); err != nil {
		return err
	}

	return nil
}

func DropAll(db *gorm.DB) error {
	return db.Migrator().DropTable(
		"missions_agents",
		&Task{},
		&Mission{},
		&Agent{},
	)
}
