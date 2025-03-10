package alice

import (
	"fmt"
	"github.com/habiliai/alice/api/aliceapi"
	"github.com/habiliai/alice/api/config"
	"github.com/habiliai/alice/api/internal/di"
	interceptors "github.com/habiliai/alice/api/internal/grpc-interceptors"
	"github.com/habiliai/alice/api/internal/mylog"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
)

func newServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			ctx = di.WithContainer(ctx, di.EnvProd)

			logger := di.MustGet[*mylog.Logger](ctx, mylog.Key)
			conf := di.MustGet[config.AliceConfig](ctx, config.AliceConfigKey)
			aliceApiServer := di.MustGet[aliceapi.AliceApiServer](ctx, aliceapi.ServerKey)

			logger.Debug("start server", "config", conf)

			grpcServer := grpc.NewServer(
				grpc.UnaryInterceptor(interceptors.NewUnaryServerInterceptor(ctx)),
			)
			grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
			aliceapi.RegisterAliceApiServer(grpcServer, aliceApiServer)

			eg := errgroup.Group{}
			eg.Go(func() error {
				address := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
				listener, err := new(net.ListenConfig).Listen(ctx, "tcp", address)
				if err != nil {
					return errors.WithStack(err)
				}
				defer listener.Close()

				logger.Info("serving grpc", "address", address)
				go func() {
					<-ctx.Done()
					grpcServer.GracefulStop()
					logger.Info("grpc server stopped")
				}()
				return errors.WithStack(grpcServer.Serve(listener))
			})

			grpcWebServer := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
				return true
			}))

			cors := func() *cors.Cors {
				if conf.IncludeDebug {
					return cors.AllowAll()
				} else {
					return cors.New(cors.Options{
						AllowedOrigins: []string{
							"https://habili.ai",
						},
						AllowedMethods: []string{
							http.MethodHead,
							http.MethodGet,
							http.MethodPost,
							http.MethodPut,
							http.MethodPatch,
							http.MethodDelete,
						},
						AllowedHeaders:   []string{"*"},
						AllowCredentials: false,
					})
				}
			}()

			eg.Go(func() error {
				httpServer := http.Server{Handler: cors.Handler(grpcWebServer)}
				address := fmt.Sprintf("%s:%d", conf.Host, conf.WebPort)
				var lc net.ListenConfig
				listener, err := lc.Listen(ctx, "tcp", address)
				if err != nil {
					return errors.WithStack(err)
				}
				defer listener.Close()

				logger.Info("serving grpc web", "address", address)

				go func() {
					<-ctx.Done()
					httpServer.Close()
					logger.Info("http server stopped")
				}()
				if err := httpServer.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
					return errors.WithStack(err)
				}
				return nil
			})

			return eg.Wait()
		},
	}

	return cmd
}
