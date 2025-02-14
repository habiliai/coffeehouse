package action

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/digo"
	hablog "github.com/habiliai/habiliai/api/pkg/log"
)

type Service interface {
	Dispatch(ctx context.Context, actionName string, args []byte) ([]byte, error)
}

type service struct {
}

const (
	ServiceKey digo.ObjectKey = "action"
)

var (
	logger = hablog.GetLogger()
)

func init() {
	digo.Register(ServiceKey, func(container *digo.Container) (any, error) {
		return &service{}, nil
	})
}
