package habgrpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

func WaitForServing(
	ctx context.Context,
	addr string,
) error {
	for interrupted := false; !interrupted; {
		select {
		case <-ctx.Done():
			interrupted = true
		case <-time.After(250 * time.Millisecond):
			// Perform gRPC health check
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return errors.Wrapf(err, "failed to connect to gRPC server")
			}
			defer conn.Close()

			healthClient := grpc_health_v1.NewHealthClient(conn)
			if resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{}); err != nil {
				return errors.Wrapf(err, "health check failed")
			} else if resp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
				interrupted = true
				break
			} else {
				return errors.Errorf("gRPC server not in serving state: %v", resp.Status)
			}
		}
	}

	return nil
}
