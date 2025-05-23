package domain

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model

	Name        string `gorm:"index:idx_agent_name,unique,where=deleted_at IS NULL"`
	AssistantId string
	IconUrl     string
	Role        string

	IncludeQuestionIntent bool
}

func (a *Agent) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(a).Error, "failed to save agent")
}
