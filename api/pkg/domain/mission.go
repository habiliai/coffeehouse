package domain

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Mission struct {
	gorm.Model

	Name string

	Agents []Agent `gorm:"many2many:missions_agents;"`
	Tasks  []Task  `gorm:"foreignKey:MissionID"`
}

func (m *Mission) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(m).Error, "failed to save mission")
}
