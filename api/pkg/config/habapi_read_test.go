package config_test

import (
	"github.com/habiliai/alice/api/pkg/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestReadHabApiConfig(t *testing.T) {
	os.Setenv("DB_HOST", "postgres.local")
	cfg := config.ReadHabApiConfig("")

	require.Equal(t, cfg.DB.Host, "postgres.local")
}
