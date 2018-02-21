package vhost

import "github.com/spf13/cobra"

func VHostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vhost",
		Short: "manage virtual hosts",
	}
	pflags := cmd.PersistentFlags()
	var output string
	pflags.StringVarP(&output, "output", "o", "", "output format yaml|json")
	cmd.AddCommand(createCmd(), deleteCmd(), getCmd(), updateCmd())
	return cmd
}
