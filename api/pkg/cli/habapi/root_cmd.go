package habapi

import "github.com/spf13/cobra"

func (c *cli) newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "habapi",
		Short: "Habapi is a tool for managing Hab API",
	}

	cmd.AddCommand(
		c.newServeCmd(),
	)

	return cmd
}
