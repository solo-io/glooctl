package vhost

import (
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/glooctl/pkg/virtualhost"
	"github.com/spf13/cobra"
)

var (
	cliOpts = &virtualhost.Options{}
)

// Cmd command to manage virtual hosts
func Cmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "virtualhost",
		Short: "manage virtual hosts",
	}
	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&cliOpts.Output, "output", "o", "", "output format yaml|json|template")
	pflags.StringVarP(&cliOpts.Template, "template", "t", "", "output template")
	cmd.AddCommand(createCmd(opts), deleteCmd(opts), getCmd(opts),
		updateCmd(opts), editCmd(opts))
	return cmd
}
