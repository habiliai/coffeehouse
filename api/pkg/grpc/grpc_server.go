package afgrpc

import (
	"context"
	"fmt"
	"github.com/habiliai/habiliai/api/pkg/digo"
	aferrors "github.com/habiliai/habiliai/api/pkg/errors"
	"github.com/habiliai/habiliai/api/pkg/habapi"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/habiliai/habiliai/api/pkg/services"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
)

const (
	ServerKey digo.ObjectKey = "proto.grpcServer"
)

func handleErrorToGrpcStatus(err error) error {
	if err == nil {
		return nil
	}

	logger.Error(fmt.Sprintf("grpc route error: %+v", err))
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, aferrors.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, aferrors.ErrBadRequest) {
		return status.Error(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, aferrors.ErrForbidden) {
		return status.Error(codes.PermissionDenied, err.Error())
	} else if errors.Is(err, aferrors.ErrUnauthorized) {
		return status.Error(codes.Unauthenticated, err.Error())
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return status.Error(codes.AlreadyExists, err.Error())
	} else if errors.Is(err, aferrors.ErrPreconditionFailed) {
		return status.Error(codes.FailedPrecondition, err.Error())
	} else if errors.Is(err, aferrors.ErrPreconditionRequired) {
		return status.Error(codes.FailedPrecondition, err.Error())
	} else {
		return status.Error(codes.Internal, err.Error())
	}
}

func createGrpcServer(
	db *gorm.DB,
	afbServer afb.AgentFatherBackendServer,
) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(func(
			ctx context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (resp interface{}, rErr error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				md = metadata.New(nil)
			}

			if values, ok := md["x-device-id"]; ok && len(values) > 0 {
				deviceId := strings.TrimSpace(values[0])
				ctx = helpers.WithDeviceId(ctx, deviceId)
				logger.Debug("metadata", "x-device-id", deviceId)
			}

			if values, ok := md["authorization"]; ok && len(values) > 0 {
				token, ok := strings.CutPrefix(values[0], "Bearer")
				if !ok {
					return nil, errors.Wrapf(aferrors.ErrUnauthorized, "Invalid authorization header")
				}
				token = strings.TrimSpace(token)
				ctx = helpers.WithAuthToken(ctx, token)
				logger.Debug("metadata", "token", token)
			}

			if values, ok := md["x-github-token"]; ok && len(values) > 0 {
				token := strings.TrimSpace(values[0])
				ctx = helpers.WithGithubToken(ctx, token)
				logger.Debug("metadata", "x-github-token", token)
			}

			rErr = db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
				ctx = helpers.WithTx(ctx, tx)

				logger.Info("call", "method", info.FullMethod)
				resp, err = handler(ctx, req)
				return
			})

			rErr = handleErrorToGrpcStatus(rErr)

			return
		}),
	)
	afb.RegisterAgentFatherBackendServer(grpcServer, afbServer)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	return grpcServer, nil
}

func init() {
	digo.Register(ServerKey, func(ctx *digo.Container) (interface{}, error) {
		db, err := digo.Get[*gorm.DB](ctx, services.ServiceKeyDB)
		if err != nil {
			return nil, err
		}

		afbServer, err := digo.Get[afb.UnsafeAgentFatherBackendServer](ctx, afb.ServerKey)
		if err != nil {
			return nil, err
		}

		return createGrpcServer(db, afbServer)
	})
}
