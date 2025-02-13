package config

import (
	"fmt"
	"github.com/habiliai/habiliai/api/pkg/constants"
)

type (
	DBConfig struct {
		PingTimeout     string
		AutoMigration   bool
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime string
		Host            string
		Port            int
		User            string
		Name            string
		Password        string
	}
)

func (c DBConfig) GetURI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, constants.SchemaName,
	)
}
