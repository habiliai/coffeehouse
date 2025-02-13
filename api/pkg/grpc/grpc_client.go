package habgrpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func NewClient(
	address string,
	ssl bool,
	timeout time.Duration,
) (*grpc.ClientConn, error) {
	var (
		creds credentials.TransportCredentials
		err   error
	)
	if !ssl {
		creds = insecure.NewCredentials()
	} else {
		creds, err = credentials.NewClientTLSFromFile("", "")
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create client tls")
		}
	}

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx, cancel := context.WithTimeoutCause(ctx, timeout, errors.Wrapf(context.DeadlineExceeded, "grpc '%s' is timed out", method))
			defer cancel()

			return errors.WithStack(invoker(ctx, method, req, reply, cc, opts...))
		}),
	)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial")
	}

	return conn, nil
}
