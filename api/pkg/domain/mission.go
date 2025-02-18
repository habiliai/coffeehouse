package domain

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Mission struct {
	gorm.Model

	Name string `gorm:"index:idx_mission_name,unique,where=deleted_at IS NULL"`

	Steps []Step `gorm:"foreignKey:MissionID"`
}

func (m *Mission) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(m).Error, "failed to save mission")
}
