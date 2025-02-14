package habapi

import (
	"github.com/habiliai/habiliai/api/pkg/action"
	"github.com/habiliai/habiliai/api/pkg/digo"
	hablog "github.com/habiliai/habiliai/api/pkg/log"
	"github.com/openai/openai-go"
)

type server struct {
	UnsafeHabiliApiServer

	openai        *openai.Client
	actionService action.Service
}

var (
	ServerKey digo.ObjectKey  = "afb.server"
	_         HabiliApiServer = (*server)(nil)
	logger                    = hablog.GetLogger()
)

func init() {
	digo.Register(ServerKey, func(container *digo.Container) (any, error) {
		openaiClient := openai.NewClient()

		actionService, err := digo.Get[action.Service](container, action.ServiceKey)
		if err != nil {
			return nil, err
		}

		return &server{
			openai:        openaiClient,
			actionService: actionService,
		}, nil
	})
}
