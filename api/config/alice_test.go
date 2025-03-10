package config_test

import (
	"github.com/habiliai/alice/api/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestReadAliceApiConfig(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres.local")
	cfg, err := config.ResolveAliceConfig("")
	require.NoError(t, err)

	require.Contains(t, cfg.DatabaseUrl, "postgres.local")
}
