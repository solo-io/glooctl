package upstream

import (
	"fmt"
	"os"

	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/spf13/cobra"
)

func deleteCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "delete upstream",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				os.Exit(1)
			}

			if err := runDelete(sc, args[0]); err != nil {
				fmt.Printf("Unable to delete upstream %s: %q\n", args[0], err)
				os.Exit(1)
			}
			fmt.Printf("Upstream %s deleted\n", args[0])
		},
	}
	return cmd
}

func runDelete(sc storage.Interface, name string) error {
	if name == "" {
		return fmt.Errorf("missing name of upstream to delete")
	}

	return sc.V1().Upstreams().Delete(name)
}
