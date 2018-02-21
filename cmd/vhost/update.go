package vhost

import (
	"fmt"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			err = runUpdate(sc, filename)
			if err != nil {
				fmt.Printf("Unable to update virtual host %q\n", err)
				return
			}
			fmt.Println("Virtual host updated")
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "file to use to update virtual host")
	cmd.MarkFlagFilename("filename")
	cmd.MarkFlagRequired("filename")
	return cmd
}

func runUpdate(sc storage.Interface, filename string) error {
	return fmt.Errorf("not implemented")
}
