package domain

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model

	MissionID uint
	Mission   Mission `gorm:"foreignKey:MissionID"`

	SeqNo int
}

func (t *Task) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(t).Error, "failed to save task")
}
