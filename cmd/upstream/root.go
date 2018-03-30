package upstream

import (
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func UpstreamCmd(opts *util.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upstream",
		Short: "manage upstreams",
	}
	pflags := cmd.PersistentFlags()
	var output string
	pflags.StringVarP(&output, "output", "o", "", "output format yaml|json")
	cmd.AddCommand(createCmd(opts), deleteCmd(opts), getCmd(opts), updateCmd(opts),
		editCmd(opts))
	return cmd
}
