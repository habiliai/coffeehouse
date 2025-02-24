package habapi

import (
	"context"
	"fmt"
	"github.com/habiliai/alice/api/pkg/callbacks"
	"github.com/habiliai/alice/api/pkg/digo"
	hablog "github.com/habiliai/alice/api/pkg/log"
	"github.com/habiliai/alice/api/pkg/services"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"runtime"
	"sync"
)

type server struct {
	UnsafeHabiliApiServer

	openai        *openai.Client
	actionService callbacks.Service

	runWg    sync.WaitGroup
	runReqCh chan runRequest
	db       *gorm.DB
}

var (
	ServerKey digo.ObjectKey = "afb.server"
	logger                   = hablog.GetLogger()
)

func (s *server) Close() {
	close(s.runReqCh)
	s.runWg.Wait()
	logger.Info("server workers all stopped")
}

func newServer(ctx context.Context, openai *openai.Client, actionService callbacks.Service, db *gorm.DB) *server {
	s := &server{
		db:            db,
		openai:        openai,
		actionService: actionService,
		runReqCh:      make(chan runRequest),
	}

	for i := 0; i < min(runtime.GOMAXPROCS(0), 8); i++ {
		workerName := fmt.Sprintf("worker-%d", i)
		s.runWg.Add(1)
		go func(workerName string) {
			defer s.runWg.Done()
			s.doRunner(ctx, workerName)
		}(workerName)
	}

	return s
}

func init() {
	digo.Register(ServerKey, func(container *digo.Container) (any, error) {
		actionService, err := digo.Get[callbacks.Service](container, callbacks.ServiceKey)
		if err != nil {
			return nil, err
		}

		db, err := digo.Get[*gorm.DB](container, services.ServiceKeyDB)
		if err != nil {
			return nil, err
		}

		var server *server
		switch container.Env {
		case digo.EnvProd:
			openaiClient := openai.NewClient(
				option.WithAPIKey(container.Config.OpenAI.ApiKey),
			)
			server = newServer(container.Context, openaiClient, actionService, db)
		case digo.EnvTest:
			openaiClient := openai.NewClient()
			server = newServer(container.Context, openaiClient, actionService, db)
		default:
			return nil, errors.Errorf("unknown environment '%s'", container.Env)
		}

		go func() {
			<-container.Context.Done()
			server.Close()
		}()

		return server, nil
	})
}
