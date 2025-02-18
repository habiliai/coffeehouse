package domain

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Action struct {
	gorm.Model

	StepID uint
	Step   Step `gorm:"foreignKey:StepID"`

	AgentID uint
	Agent   Agent `gorm:"foreignKey:AgentID"`

	Subject string
}

func (a *Action) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(a).Error, "failed to save action")
}
