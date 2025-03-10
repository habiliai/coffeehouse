package interceptors

import (
	"context"
	myerrors "github.com/habiliai/agentruntime/errors"
	"github.com/habiliai/alice/api/internal/di"
	"github.com/habiliai/alice/api/internal/mylog"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	logger := di.MustGet[*mylog.Logger](ctx, mylog.Key)

	logger.Error("error in gRPC handler", "err", err)
	if errors.Is(err, myerrors.ErrNotFound) {
		return status.Errorf(codes.NotFound, "%v", err)
	} else if errors.Is(err, myerrors.ErrInvalidParams) {
		return status.Errorf(codes.InvalidArgument, "%v", err)
	} else {
		return status.Errorf(codes.Internal, "%v", err)
	}
}
