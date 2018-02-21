package vhost

import (
	"fmt"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			err = runCreate(sc, filename)
			if err != nil {
				fmt.Printf("Unable to create virtual host %q\n", err)
				return
			}
			fmt.Println("Virtual host created")
		},
	}

	cmd.Flags().StringVar(&filename, "filename", "f", "file to use to create virtual host")
	return cmd
}

func runCreate(sc storage.Interface, filename string) error {
	return fmt.Errorf("not implemented")
}
