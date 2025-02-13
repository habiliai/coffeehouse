package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/config"
	hablog "github.com/habiliai/habiliai/api/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type cli struct {
	cfg config.HabApiConfig
}

func (c *cli) ReadInConfig() error {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err == nil {
		logger.Info("Read in config from .env")
	}

	if err := v.Unmarshal(&c.cfg); err != nil {
		return errors.Wrapf(err, "failed to unmarshal config")
	}

	return nil
}

var (
	logger = hablog.GetLogger()
)

func Execute(ctx context.Context) error {
	cli := &cli{}
	return cli.newRootCmd().ExecuteContext(ctx)
}
