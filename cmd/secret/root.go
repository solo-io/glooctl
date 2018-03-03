package secret

import "github.com/spf13/cobra"

const (
	flagFilename = "filename"
)

func SecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "manage secrets for upstreams in gloo",
	}
	cmd.AddCommand(createCmd(), deleteCmd())
	return cmd
}
