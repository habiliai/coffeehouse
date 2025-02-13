package domain

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE SCHEMA IF NOT EXISTS afb`).Error; err != nil {
		return errors.Wrapf(err, "failed to create schema")
	}

	if err := errors.Wrapf(db.
		AutoMigrate(), "failed to auto migrate"); err != nil {
		return err
	}

	return nil
}

func DropAll(db *gorm.DB) error {
	return db.Migrator().DropTable()
}
