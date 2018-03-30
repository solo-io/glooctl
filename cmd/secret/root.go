package secret

import (
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

const (
	flagFilename = "filename"
)

func SecretCmd(opts *util.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "manage secrets for upstreams in gloo",
	}
	cmd.AddCommand(createCmd(opts), deleteCmd(opts))
	return cmd
}
