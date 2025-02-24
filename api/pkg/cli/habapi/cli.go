package habapi

import (
	"context"
	habconfig "github.com/habiliai/alice/api/pkg/config"
	hablog "github.com/habiliai/alice/api/pkg/log"
)

type cli struct {
	cfg habconfig.HabApiConfig
}

func (c *cli) ReadInConfig() error {
	// Set default values
	c.cfg = habconfig.ReadHabApiConfig("")
	return nil
}

var (
	logger = hablog.GetLogger()
)

func Execute(ctx context.Context) error {
	cli := &cli{}
	return cli.newRootCmd().ExecuteContext(ctx)
}
