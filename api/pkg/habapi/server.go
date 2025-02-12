package habapi

import "github.com/habiliai/habiliai/api/pkg/digo"

type server struct {
	UnsafeAgentFatherBackendServer
}

var (
	ServerKey digo.ObjectKey           = "afb.server"
	_         AgentFatherBackendServer = (*server)(nil)
)

func init() {
	digo.Register(ServerKey, func(container *digo.Container) (any, error) {
		return &server{}, nil
	})
}
