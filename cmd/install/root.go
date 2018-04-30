package install

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install gloo on different platforms",
	}
	cmd.AddCommand(dockerCmd())
	return cmd
}
