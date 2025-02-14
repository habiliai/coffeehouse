package habapi

import (
	"context"
	"github.com/golobby/dotenv"
	"github.com/golobby/env/v2"
	habconfig "github.com/habiliai/habiliai/api/pkg/config"
	hablog "github.com/habiliai/habiliai/api/pkg/log"
	"os"
)

type cli struct {
	cfg habconfig.HabApiConfig
}

func (c *cli) ReadInConfig() error {
	// Set default values
	c.cfg = habconfig.HabApiConfig{
		Address:      "",
		Port:         8000,
		WebPort:      8001,
		IncludeDebug: true,
		DB: habconfig.DBConfig{
			PingTimeout:     "5s",
			AutoMigration:   true,
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: "1h",
			Host:            "localhost",
			Port:            5432,
			User:            "habiliai",
			Name:            "habiliai",
			Password:        "habiliai",
		},
		OpenAI: habconfig.OpenAIConfig{
			ApiKey: "",
		},
	}

	envFile, err := os.Open(".env")
	defer func() {
		if envFile != nil {
			envFile.Close()
		}
	}()
	if err != nil {
		logger.Warn("failed to open .env file", "error", err)
	} else if err := dotenv.NewDecoder(envFile).Decode(&c.cfg); err != nil {
		logger.Warn("failed to decode .env file", "error", err)
	} else {
		logger.Info("Read in .env", "config", c.cfg)
	}

	if err := env.Feed(&c.cfg); err == nil {
		logger.Info("Read in env", "config", c.cfg)
	} else {
		logger.Warn("failed to read in env", "error", err)
	}

	logger.Info("Read in config", "config", c.cfg)

	return nil
}

var (
	logger = hablog.GetLogger()
)

func Execute(ctx context.Context) error {
	cli := &cli{}
	return cli.newRootCmd().ExecuteContext(ctx)
}
