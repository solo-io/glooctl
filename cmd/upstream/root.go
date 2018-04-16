package upstream

import (
	"github.com/solo-io/glooctl/pkg/client"
	"github.com/solo-io/glooctl/pkg/upstream"
	"github.com/spf13/cobra"
)

var (
	cliOpts = &upstream.Options{}
)

func UpstreamCmd(opts *client.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upstream",
		Short: "manage upstreams",
	}
	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&cliOpts.Output, "output", "o", "", "output format yaml|json|template")
	pflags.StringVarP(&cliOpts.Template, "template", "t", "", "output template")
	cmd.AddCommand(createCmd(opts), deleteCmd(opts), getCmd(opts), updateCmd(opts),
		editCmd(opts))
	return cmd
}
