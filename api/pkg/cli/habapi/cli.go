package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/config"
	hablog "github.com/habiliai/habiliai/api/pkg/log"
)

type cli struct {
	cfg config.HabApiConfig
}

var (
	logger = hablog.GetLogger()
)

func Execute(ctx context.Context) error {
	cli := &cli{}
	return cli.newRootCmd().ExecuteContext(ctx)
}
