package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list configuration",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("config list - not implemented")
		},
	}
	return cmd
}
