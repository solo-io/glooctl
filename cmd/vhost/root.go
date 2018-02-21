package vhost

import "github.com/spf13/cobra"

func VHostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vhost",
		Short: "manage virtual hosts",
	}
	cmd.AddCommand(createCmd(), deleteCmd(), getCmd(), updateCmd())
	return cmd
}
