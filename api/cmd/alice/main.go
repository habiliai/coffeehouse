package main

import (
	"context"
	"fmt"
	"github.com/habiliai/alice/api/cli/alice"
	_ "go.uber.org/automaxprocs"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	defer cancel()

	cmd := alice.NewRootCmd()
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "exit by error: %+v", err)
		os.Exit(1)
	}
}
