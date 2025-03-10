package interceptors

import (
	"context"
	"github.com/habiliai/alice/api/internal/di"
	"github.com/habiliai/alice/api/internal/mylog"
	"google.golang.org/grpc"
)

func NewUnaryServerInterceptor(ctx context.Context) func(context.Context, any, *grpc.UnaryServerInfo, grpc.UnaryHandler) (resp any, err error) {
	logger := di.MustGet[*mylog.Logger](ctx, mylog.Key)
	return func(_ context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		logger.Info("[gRPC] call", "path", info.FullMethod)

		resp, err = handler(ctx, req)

		err = handleError(ctx, err)
		return
	}
}
