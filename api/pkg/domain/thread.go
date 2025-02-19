package domain

import (
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Thread struct {
	gorm.Model

	OpenaiThreadId string
	CurrentRunId   string
	LastMessageId  string

	MissionID uint
	Mission   Mission `gorm:"foreignKey:MissionID"`

	CurrentStepSeqNo int
	AllDone          bool
	Result           string

	Data datatypes.JSONType[map[string]any]
}

func (t *Thread) Save(db *gorm.DB) error {
	return errors.Wrapf(db.Save(t).Error, "failed to save thread")
}

func (t *Thread) GetCurrentStep() (*Step, error) {
	for _, step := range t.Mission.Steps {
		if step.SeqNo == t.CurrentStepSeqNo {
			return &step, nil
		}
	}
	return nil, errors.Errorf("failed to find current step with seq no %d", t.CurrentStepSeqNo)
}
