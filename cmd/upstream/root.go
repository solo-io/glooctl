package upstream

import "github.com/spf13/cobra"

func UpstreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upstream",
		Short: "manage upstreams",
	}
	pflags := cmd.PersistentFlags()
	var output string
	pflags.StringVarP(&output, "output", "o", "", "output format yaml|json")
	cmd.AddCommand(createCmd(), deleteCmd(), getCmd(), updateCmd())
	return cmd
}
