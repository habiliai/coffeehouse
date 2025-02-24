package domain

import (
	"fmt"
	"github.com/habiliai/alice/api/pkg/constants"
	hablog "github.com/habiliai/alice/api/pkg/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	logger = hablog.GetLogger()
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, constants.SchemaName)).Error; err != nil {
		return errors.Wrapf(err, "failed to create schema")
	}

	if err := errors.Wrapf(db.
		AutoMigrate(
			&Agent{},
			&Mission{},
			&Thread{},
			&Step{},
			&Action{},
			&ActionWork{},
			&AgentWork{},
		), "failed to auto migrate"); err != nil {
		return err
	}

	return nil
}

func DropAll(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&AgentWork{},
		&ActionWork{},
		&Action{},
		&Step{},
		&Thread{},
		&Mission{},
		&Agent{},
	)
}
