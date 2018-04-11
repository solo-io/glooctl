package secret

import (
	"github.com/solo-io/glooctl/pkg/client"
	"github.com/spf13/cobra"
)

const (
	flagFilename = "filename"
)

func SecretCmd(opts *client.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "manage secrets for upstreams in gloo",
	}
	cmd.AddCommand(createCmd(opts), deleteCmd(opts), getCmd(opts))
	return cmd
}
