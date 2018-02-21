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
		Use:   "delete",
		Short: "delete virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			err = runDelete(sc, filename)
			if err != nil {
				fmt.Printf("Unable to delete virtual host %q\n", err)
				return
			}
			fmt.Println("Virtual host deleted")
		},
	}

	cmd.Flags().StringVar(&filename, "filename", "f", "file to use to delete virtual host")
	return cmd
}

func runDelete(sc storage.Interface, filename string) error {
	return fmt.Errorf("not implemented")
}
