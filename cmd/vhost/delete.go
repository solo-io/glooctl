package vhost

import (
	"fmt"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func deleteCmd() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "delete virtual host",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			sc, errj := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			name := args[0]
			if err := runDelete(sc, name); err != nil {
				fmt.Printf("Unable to delete virtual host %s: %q\n", name, err)
				return
			}
			fmt.Printf("Virtual host %s deleted\n", name)
		},
	}
	return cmd
}

func runDelete(sc storage.Interface, name string) error {
	if name == "" {
		return fmt.Errorf("missing name of virtual host to delete")
	}
	return sc.V1().VirtualHosts().Delete(name)
}
