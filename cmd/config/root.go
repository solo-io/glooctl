package config

import "github.com/spf13/cobra"

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "configuration for glooctl",
	}

	cmd.AddCommand(listCmd(), setCmd())
	return cmd
}
