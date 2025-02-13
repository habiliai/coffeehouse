package habapi

import (
	"github.com/habiliai/habiliai/api/pkg/digo"
	hablog "github.com/habiliai/habiliai/api/pkg/log"
	"github.com/openai/openai-go"
)

type server struct {
	UnsafeHabiliApiServer

	openai *openai.Client
}

var (
	ServerKey digo.ObjectKey  = "afb.server"
	_         HabiliApiServer = (*server)(nil)
	logger                    = hablog.GetLogger()
)

func init() {
	digo.Register(ServerKey, func(container *digo.Container) (any, error) {
		return &server{}, nil
	})
}
