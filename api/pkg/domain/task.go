package domain

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model

	MissionID uint    `gorm:"index:idx_mission_id_seqno,unique,where=deleted_at IS NULL"`
	Mission   Mission `gorm:"foreignKey:MissionID"`

	SeqNo int `gorm:"index:idx_mission_id_seqno,unique,where=deleted_at IS NULL"`
}

func (t *Task) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(t).Error, "failed to save task")
}
