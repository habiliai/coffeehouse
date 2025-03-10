package aliceapi

import (
	"context"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/agentruntime/runtime"
	"github.com/habiliai/agentruntime/thread"
	"github.com/habiliai/alice/api/internal/agentruntimeclient"
	"github.com/habiliai/alice/api/internal/db"
	"github.com/habiliai/alice/api/internal/di"
	"github.com/habiliai/alice/api/internal/mylog"
	"gorm.io/gorm"
	"sync"
)

type server struct {
	UnsafeAliceApiServer

	db *gorm.DB
	wg sync.WaitGroup

	threadManager thread.ThreadManagerClient
	agentManager  agent.AgentManagerClient
	runtime       runtime.AgentRuntimeClient
	logger        *mylog.Logger
}

var (
	ServerKey = di.NewKey()
)

func (s *server) Close() error {
	s.logger.Info("server workers all stopped")
	s.wg.Wait()
	return nil
}

func init() {
	di.Register(ServerKey, func(ctx context.Context, env di.Env) (any, error) {
		db, err := di.Get[*gorm.DB](ctx, db.Key)
		if err != nil {
			return nil, err
		}

		threadManager := agentruntimeclient.GetThreadManager(ctx)
		agentManager := agentruntimeclient.GetAgentManager(ctx)
		runtimeClient := agentruntimeclient.GetAgentRuntime(ctx)

		s := &server{
			db:            db,
			threadManager: threadManager,
			agentManager:  agentManager,
			runtime:       runtimeClient,
			logger:        di.MustGet[*mylog.Logger](ctx, mylog.Key),
		}

		go func() {
			<-ctx.Done()
			s.Close()
		}()

		return s, nil
	})
}
