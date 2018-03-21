package secret

import (
	"fmt"

	secret "github.com/solo-io/gloo-secret"
	"github.com/solo-io/glooctl/pkg/secrets"
	"github.com/spf13/cobra"
)

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "delete secret",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			si, err := secrets.GetSecretClient(c)
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

func runDelete(si secret.SecretInterface, name string) error {
	if name == "" {
		return fmt.Errorf("missing name of secret to delete")
	}

	return si.V1().Delete(name)
}
