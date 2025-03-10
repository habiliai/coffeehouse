package domain

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ActionWork struct {
	ActionID uint   `gorm:"primarykey"`
	Action   Action `gorm:"foreignKey:ActionID"`
	ThreadID uint   `gorm:"primarykey"`
	Thread   Thread `gorm:"foreignKey:ThreadID"`

	Done  bool
	Error string
}

func (aw *ActionWork) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(aw).Error, "failed to save action work")
}
