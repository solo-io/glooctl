package route

import "github.com/spf13/cobra"

// RouteCmd returns command related to managing routes on a virtual host
func RouteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "manage routes on a virtual host",
	}

	pflags := cmd.PersistentFlags()
	var output string
	pflags.StringVarP(&output, "output", "o", "", "output format yaml|json")
	var vhost string
	pflags.StringVarP(&vhost, "vhost", "V", "", "name of the virtual host")
	cmd.MarkPersistentFlagRequired("vhost")
	cmd.AddCommand(appendCmd(), sortCmd(), getCmd())
	return cmd
}
