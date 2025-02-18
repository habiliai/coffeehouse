package callbacks

import (
	"context"
	"encoding/json"
	"github.com/habiliai/habiliai/api/pkg/config"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/pkg/errors"
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

	memory := map[string]any{}
	if len(c.Thread.Memory) > 0 {
		if err := json.Unmarshal(c.Thread.Memory, &memory); err != nil {
			return errors.Wrapf(err, "failed to unmarshal memory")
		}
	}
	for key, value := range records {
		memory[key] = value
	}

	memoryBytes, err := json.Marshal(memory)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal memory")
	}

	c.Thread.Memory = memoryBytes
	return c.Thread.Save(tx)
}
