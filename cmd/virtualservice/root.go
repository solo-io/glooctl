package virtualservice

import (
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/glooctl/pkg/virtualservice"
	"github.com/spf13/cobra"
)

var (
	cliOpts = &virtualservice.Options{}
)

// Cmd command to manage virtual services
func Cmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "virtualservice",
		Short: "manage virtual services",
	}
	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&cliOpts.Output, "output", "o", "", "output format yaml|json|template")
	pflags.StringVarP(&cliOpts.Template, "template", "t", "", "output template")
	cmd.AddCommand(createCmd(opts), deleteCmd(opts), getCmd(opts),
		updateCmd(opts), editCmd(opts))
	return cmd
}
