package secret

import (
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a secret for upstreams",
	}
	flags := cmd.PersistentFlags()
	var name string
	flags.StringVar(&name, "name", "", "name for secret")
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(createAWS(), createGCF())
	return cmd
}
