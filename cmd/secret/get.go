package secret

import (
	"fmt"
	"os"

	"github.com/solo-io/glooctl/pkg/client"
	"github.com/solo-io/glooctl/pkg/secret"
	"github.com/spf13/cobra"
)

func getCmd(opts *client.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name (optional)]",
		Short: "get a secret or secret list",
		Args:  cobra.MaximumNArgs(1),
		Run: func(c *cobra.Command, a []string) {
			sc, err := client.StorageClient(opts)
			if err != nil {
				fmt.Println("Unable to get storage client:", err)
				os.Exit(1)
			}

			var name string
			if len(a) > 0 {
				name = a[0]
			}
			si, err := client.SecretClient(opts)
			if err != nil {
				fmt.Println("Unable to get secret client:", err)
				os.Exit(1)
			}
			err = secret.Get(sc, si, name)
			if err != nil {
				fmt.Println("Unable to get secret:", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}
