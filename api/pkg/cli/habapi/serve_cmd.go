package habapi

import (
	"fmt"
	"github.com/habiliai/alice/api/pkg/digo"
	habgrpc "github.com/habiliai/alice/api/pkg/grpc"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os/signal"
	"syscall"
)

func (c *cli) newServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
			defer cancel()

			if err := c.ReadInConfig(); err != nil {
				return err
			}

			logger.Debug("start server", "config", c.cfg)

			container := digo.NewContainer(ctx, digo.EnvProd, &c.cfg)
			grpcServer, err := digo.Get[*grpc.Server](container, habgrpc.ServerKey)
			if err != nil {
				return err
			}
			defer grpcServer.GracefulStop()
			go func() {
				<-ctx.Done()
				grpcServer.GracefulStop()
				logger.Info("grpc server stopped")
			}()

			eg := errgroup.Group{}
			eg.Go(func() error {
				address := fmt.Sprintf("%s:%d", c.cfg.Address, c.cfg.Port)
				listener, err := new(net.ListenConfig).Listen(ctx, "tcp", address)
				if err != nil {
					return errors.WithStack(err)
				}
				defer listener.Close()

				logger.Info("serving grpc", "address", address)
				return errors.WithStack(grpcServer.Serve(listener))
			})

			grpcWebServer := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
				return true
			}))

			cors := func() *cors.Cors {
				if c.cfg.IncludeDebug {
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

			httpServer := http.Server{Handler: cors.Handler(grpcWebServer)}
			defer httpServer.Close()
			go func() {
				<-ctx.Done()
				httpServer.Close()
				logger.Info("http server stopped")
			}()

			eg.Go(func() error {
				address := fmt.Sprintf("%s:%d", c.cfg.Address, c.cfg.WebPort)
				listener, err := new(net.ListenConfig).Listen(ctx, "tcp", address)
				if err != nil {
					return errors.WithStack(err)
				}
				defer listener.Close()

				logger.Info("serving grpc web", "address", address)

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
