package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func setCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "set a config value",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("set value - not implemented")
		},
	}
	return cmd
}
