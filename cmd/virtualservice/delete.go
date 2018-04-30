package virtualservice

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/spf13/cobra"
)

func deleteCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "delete virtual service",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			name := args[0]
			if err := runDelete(sc, name); err != nil {
				fmt.Printf("Unable to delete virtual service %s: %q\n", name, err)
				return
			}
			fmt.Printf("Virtual service %s deleted\n", name)
		},
	}
	return cmd
}

func runDelete(sc storage.Interface, name string) error {
	if name == "" {
		return fmt.Errorf("missing name of virtual service to delete")
	}
	return sc.V1().VirtualServices().Delete(name)
}
