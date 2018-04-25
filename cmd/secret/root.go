package secret

import (
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/spf13/cobra"
)

// Cmd command for managing secrets in Gloo
func Cmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "manage secrets used in gloo",
	}
	cmd.AddCommand(createCmd(opts), deleteCmd(opts), getCmd(opts))
	return cmd
}
