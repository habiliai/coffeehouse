package main

import (
	"context"
	"fmt"
	"github.com/habiliai/alice/api/pkg/cli/habapi"
	_ "go.uber.org/automaxprocs"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := habapi.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "exit by error: %+v", err)
		os.Exit(1)
	}
}
