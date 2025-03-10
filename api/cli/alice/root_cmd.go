package alice

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alice",
		Short: "alice is a tool for managing Hab API",
	}

	cmd.AddCommand(
		newServeCmd(),
		newSeedCmd(),
	)

	return cmd
}
