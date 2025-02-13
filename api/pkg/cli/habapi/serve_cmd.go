package habapi

import (
	"fmt"
	"github.com/habiliai/habiliai/api/pkg/digo"
	habgrpc "github.com/habiliai/habiliai/api/pkg/grpc"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func (c *cli) newServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			viper.AddConfigPath(".")
			viper.SetConfigName(".env")
			viper.SetConfigType("env")

			//viper.AutomaticEnv()

			if err := viper.ReadInConfig(); err != nil {
				return errors.Wrapf(err, "failed to read in config")
			}

			if err := viper.Unmarshal(&c.cfg); err != nil {
				return errors.Wrapf(err, "failed to unmarshal config")
			}

			logger.Debug("start server", "config", c.cfg)

			container := digo.NewContainer(cmd.Context(), digo.EnvProd, &c.cfg)
			grpcServer, err := digo.Get[*grpc.Server](container, habgrpc.ServerKey)
			if err != nil {
				return err
			}
			defer grpcServer.GracefulStop()
			go func() {
				<-ctx.Done()
				grpcServer.GracefulStop()
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

			httpServer := http.Server{Handler: grpcWebServer}
			defer httpServer.Close()
			go func() {
				<-ctx.Done()
				httpServer.Close()
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
