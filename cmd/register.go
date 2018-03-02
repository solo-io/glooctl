package cmd

import (
	"fmt"

	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func registerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register gloo resources",
		Run: func(c *cobra.Command, args []string) {
			storageClient, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to register resource defintions %q\n", err)
				return
			}
			err = storageClient.V1().Register()
			if err != nil && !storage.IsAlreadyExists(err) {
				fmt.Printf("Unable to register resource definitions %q\n", err)
				return
			}
			fmt.Println("Registered resource definitions.")
		},
	}
	return cmd
}
