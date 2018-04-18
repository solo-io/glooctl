package secret

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/secretstorage"
	"github.com/solo-io/gloo/pkg/storage/dependencies"
	"github.com/spf13/cobra"
)

func deleteCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "delete secret",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			si, err := secretstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Println("Unable to create secret client:", err)
				return
			}

			if err := runDelete(si, args[0]); err != nil {
				fmt.Printf("Unable to delete secret %s: %q\n", args[0], err)
				return
			}
			fmt.Printf("Secret %s deleted\n", args[0])
		},
	}
	return cmd
}

func runDelete(si dependencies.SecretStorage, name string) error {
	if name == "" {
		return fmt.Errorf("missing name of secret to delete")
	}

	return si.Delete(name)
}
