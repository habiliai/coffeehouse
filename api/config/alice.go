package config

import (
	"context"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/habiliai/alice/api/internal/di"
	"github.com/pkg/errors"
	"os"
)

type AliceConfig struct {
	LogLevel   string `env:"LOG_LEVEL"`
	LogHandler string `env:"LOG_HANDLER"`

	Host         string `env:"HOST"`
	Port         int    `env:"PORT"`
	WebPort      int    `env:"WEB_PORT"`
	IncludeDebug bool   `env:"INCLUDE_DEBUG"`

	DatabaseUrl         string `env:"DATABASE_URL"`
	DatabaseAutoMigrate bool   `env:"DATABASE_AUTO_MIGRATE"`

	OpenAIApiKey string `env:"OPENAI_API_KEY"`
	LumaApiKey   string `env:"LUMA_API_KEY"`

	OpenWeatherApiKey string `env:"OPENWEATHER_API_KEY"`
	Twitter           TwitterConfig

	AgentRuntimeEndpoint string `env:"AGENT_RUNTIME_ENDPOINT"`
}

func ResolveAliceConfig(configFile string) (cfg AliceConfig, err error) {
	// Set default values
	cfg = AliceConfig{
		Host:                "0.0.0.0",
		Port:                8000,
		WebPort:             8001,
		IncludeDebug:        true,
		DatabaseUrl:         "postgres://habilaii:habiliai@localhost:5432/habiliai?sslmode=disable",
		OpenAIApiKey:        "",
		DatabaseAutoMigrate: true,
		LogLevel:            "INFO",
	}

	cfgReader := config.New()
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		cfgReader.AddFeeder(feeder.DotEnv{Path: ".env"})
	}
	if configFile != "" {
		if _, err = os.Stat(configFile); os.IsNotExist(err) {
			err = errors.Wrapf(err, "failed to find config file: %s", configFile)
			return
		}
		cfgReader.AddFeeder(feeder.DotEnv{Path: configFile})
	}
	cfgReader.AddFeeder(feeder.Env{})
	cfgReader.AddStruct(&cfg)

	if err = errors.Wrapf(cfgReader.Feed(), "failed to read config"); err != nil {
		return
	}

	return cfg, nil
}

var (
	AliceConfigKey = di.NewKey()
)

func init() {
	di.Register(AliceConfigKey, func(ctx context.Context, env di.Env) (any, error) {
		projectDir := os.Getenv("PROJECT_DIR")
		if projectDir == "" {
			projectDir = "."
		}
		
		configFile := ""
		if env == di.EnvTest {
			configFile = projectDir + "/.env.test"
		}

		return ResolveAliceConfig(configFile)
	})
}
