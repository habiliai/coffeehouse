package callbacks

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/config"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"gorm.io/datatypes"
)

type Metadata struct {
	AgentWork  *domain.AgentWork
	ActionWork *domain.ActionWork
	Thread     *domain.Thread
}

type Context struct {
	context.Context
	Metadata

	config *config.HabApiConfig
}

func (c *Context) UpdateMemory(records map[string]any) error {
	tx := helpers.GetTx(c)

	memory := c.Thread.Data.Data()
	for key, value := range records {
		memory[key] = value
	}

	c.Thread.Data = datatypes.NewJSONType(memory)
	return c.Thread.Save(tx)
}
