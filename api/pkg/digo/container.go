package digo

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/config"
	aflog "github.com/habiliai/habiliai/api/pkg/log"
)

type (
	Env       string
	Container struct {
		context.Context
		Env    Env
		Config *config.AfbConfig

		objects map[ObjectKey]any
	}
)

const (
	EnvProd = "prod"
	EnvTest = "test"
)

var (
	logger = aflog.GetLogger()
)

func NewContainer(
	ctx context.Context,
	env Env,
	cfg *config.AfbConfig,
) *Container {
	return &Container{
		Context: ctx,
		Env:     env,
		Config:  cfg,
		objects: map[ObjectKey]any{},
	}
}
