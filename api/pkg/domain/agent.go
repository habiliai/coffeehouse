package domain

import "gorm.io/gorm"

type Agent struct {
	gorm.Model

	Name        string `gorm:"index:idx_agent_name,unique,where=deleted_at IS NULL"`
	AssistantId string
	IconUrl     string
}
