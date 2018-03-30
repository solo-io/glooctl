package secret

import (
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	Name string
}

func createCmd(opts *util.StorageOptions) *cobra.Command {
	createOpts := CreateOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a secret for upstreams",
	}
	flags := cmd.PersistentFlags()
	flags.StringVar(&createOpts.Name, "name", "", "name for secret")
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(createAWS(opts, &createOpts), createGCF(opts, &createOpts))
	return cmd
}
