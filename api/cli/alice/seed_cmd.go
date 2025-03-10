package alice

import (
	"github.com/habiliai/alice/api/domain/seed"
	"github.com/habiliai/alice/api/internal/di"
	"github.com/spf13/cobra"
)

func newSeedCmd() *cobra.Command {
	flags := &struct {
		reset bool
	}{}

	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed the database",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := di.WithContainer(cmd.Context(), di.EnvProd)

			return seed.Seed(ctx, flags.reset)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&flags.reset, "reset", false, "Reset the database")

	return cmd
}
