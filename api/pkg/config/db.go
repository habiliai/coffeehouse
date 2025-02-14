package config

import (
	"fmt"
	"github.com/habiliai/habiliai/api/pkg/constants"
)

type (
	DBConfig struct {
		PingTimeout     string `env:"DB_PING_TIMEOUT"`
		AutoMigration   bool   `env:"DB_AUTO_MIGRATION"`
		MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS"`
		MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS"`
		ConnMaxLifetime string `env:"DB_CONN_MAX_LIFETIME"`
		Host            string `env:"DB_HOST"`
		Port            int    `env:"DB_PORT"`
		User            string `env:"DB_USER"`
		Name            string `env:"DB_NAME"`
		Password        string `env:"DB_PASSWORD"`
	}
)

func (c DBConfig) GetURI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, constants.SchemaName,
	)
}
