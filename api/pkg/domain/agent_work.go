package domain

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AgentWork struct {
	AgentID  uint   `gorm:"primarykey"`
	Agent    Agent  `gorm:"foreignKey:AgentID"`
	ThreadID uint   `gorm:"primarykey"`
	Thread   Thread `gorm:"foreignKey:ThreadID"`

	Status AgentStatus
}

type AgentStatus int

const (
	AgentStatusIdle AgentStatus = iota
	AgentStatusWorking
	AgentStatusWaiting
)

func (t *AgentWork) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(t).Error, "failed to save agent work")
}
