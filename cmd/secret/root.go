package secret

import (
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/spf13/cobra"
)

const (
	flagFilename = "filename"
)

func SecretCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "manage secrets for upstreams in gloo",
	}
	cmd.AddCommand(createCmd(opts), deleteCmd(opts), getCmd(opts))
	return cmd
}
