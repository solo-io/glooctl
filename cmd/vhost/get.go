package vhost

import (
	"fmt"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func getCmd() *cobra.Command {
	var output string
	var name string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			err = runGet(sc, output, name)
			if err != nil {
				fmt.Printf("Unable to create virtual host %q\n", err)
				return
			}
			fmt.Println("Virtual host created")
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "output format yaml|json")
	cmd.Flags().StringVarP(&name, "name", "n", "", "name of virtual host to get; returns all if empty")
	return cmd
}

func runGet(sc storage.Interface, output, name string) error {
	return fmt.Errorf("not implemented")
}
