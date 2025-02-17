package config

import (
	"github.com/golobby/dotenv"
	"github.com/golobby/env/v2"
	"os"
)

func ReadHabApiConfig(configFile string) HabApiConfig {
	// Set default values
	cfg := HabApiConfig{
		Address:      "",
		Port:         8000,
		WebPort:      8001,
		IncludeDebug: true,
		DB: DBConfig{
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
		OpenAI: OpenAIConfig{
			ApiKey: "",
		},
	}

	if configFile == "" {
		configFile = ".env"
	}

	envFile, err := os.Open(configFile)
	defer func() {
		if envFile != nil {
			envFile.Close()
		}
	}()
	if err != nil {
		logger.Warn("failed to open .env file", "error", err)
	} else if err := dotenv.NewDecoder(envFile).Decode(&cfg); err != nil {
		logger.Warn("failed to decode .env file", "error", err)
	} else {
		logger.Info("Read in .env", "config", cfg)
	}

	if err := env.Feed(&cfg); err == nil {
		logger.Info("Read in env", "config", cfg)
	} else {
		logger.Warn("failed to read in env", "error", err)
	}

	logger.Info("Read in config", "config", cfg)

	return cfg
}
