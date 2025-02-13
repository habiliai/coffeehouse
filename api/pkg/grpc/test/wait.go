package tclgrpctest

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"testing"
	"time"
)

func WaitForServing(
	t *testing.T,
	ctx context.Context,
	addr string,
) {
	for interrupted := false; !interrupted; {
		select {
		case <-ctx.Done():
			interrupted = true
		case <-time.After(250 * time.Millisecond):
			// Perform gRPC health check
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatalf("failed to connect to gRPC server: %v", err)
			}
			defer conn.Close()

			healthClient := grpc_health_v1.NewHealthClient(conn)
			if resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{}); err != nil {
				t.Fatalf("health check failed: %v", err)
			} else if resp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
				interrupted = true
				break
			} else {
				t.Fatalf("gRPC server not in serving state: %v", resp.Status)
			}
		}
	}
}
