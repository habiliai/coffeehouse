package callbacks

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/config"
	"github.com/habiliai/habiliai/api/pkg/digo"
	hablog "github.com/habiliai/habiliai/api/pkg/log"
)

type Service interface {
	Dispatch(ctx context.Context, actionName string, args []byte, metadata Metadata) (any, error)
}

type service struct {
	config *config.HabApiConfig
}

const (
	ServiceKey digo.ObjectKey = "action"
)

var (
	logger = hablog.GetLogger()
)

func NewService(config *config.HabApiConfig) Service {
	return &service{
		config: config,
	}
}

func init() {
	digo.Register(ServiceKey, func(container *digo.Container) (any, error) {
		return NewService(container.Config), nil
	})
}
