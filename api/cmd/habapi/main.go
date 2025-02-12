package main

import (
	"context"
	"fmt"
	"github.com/habiliai/habiliai/api/pkg/cli/habapi"
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
