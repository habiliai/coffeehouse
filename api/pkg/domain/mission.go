package domain

import (
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Mission struct {
	gorm.Model

	Name         string
	AssistantIds datatypes.JSONSlice[string]
}

func (m *Mission) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(m).Error, "failed to save mission")
}
