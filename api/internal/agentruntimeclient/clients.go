package agentruntimeclient

import (
	"context"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/agentruntime/runtime"
	"github.com/habiliai/agentruntime/thread"
	"github.com/habiliai/alice/api/config"
	"github.com/habiliai/alice/api/internal/di"
	"github.com/habiliai/alice/api/internal/mylog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	agentRuntimeClientKey = di.NewKey()
)

func GetThreadManager(ctx context.Context) thread.ThreadManagerClient {
	conn := di.MustGet[*grpc.ClientConn](ctx, agentRuntimeClientKey)
	return thread.NewThreadManagerClient(conn)
}

func GetAgentManager(ctx context.Context) agent.AgentManagerClient {
	conn := di.MustGet[*grpc.ClientConn](ctx, agentRuntimeClientKey)
	return agent.NewAgentManagerClient(conn)
}

func GetAgentRuntime(ctx context.Context) runtime.AgentRuntimeClient {
	conn := di.MustGet[*grpc.ClientConn](ctx, agentRuntimeClientKey)
	return runtime.NewAgentRuntimeClient(conn)
}

func init() {
	di.Register(agentRuntimeClientKey, func(ctx context.Context, env di.Env) (any, error) {
		conf := di.MustGet[config.AliceConfig](ctx, config.AliceConfigKey)
		logger := di.MustGet[*mylog.Logger](ctx, mylog.Key)
		conn, err := grpc.NewClient(conf.AgentRuntimeEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create grpc client")
		}
		go func() {
			<-ctx.Done()
			if err := conn.Close(); err != nil {
				logger.Warn("failed to close grpc client", "err", err)
			}
		}()

		return conn, nil
	})
}
