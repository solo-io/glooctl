package vhost

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
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
			vh, err := runCreate(sc, filename)
			if err != nil {
				fmt.Printf("Unable to create virtual host %q\n", err)
				return
			}
			fmt.Println("Virtual host created")
			output, _ := c.InheritedFlags().GetString("output")
			if output == "yaml" {
				printYAML(vh)
			}
			if output == "json" {
				printJSON(vh)
			}
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "file to use to create virtual host")
	cmd.MarkFlagFilename("filename")
	cmd.MarkFlagRequired("filename")
	return cmd
}

func runCreate(sc storage.Interface, filename string) (*v1.VirtualHost, error) {
	vh, err := parseFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load virtual host from %s", filename)
	}
	return sc.V1().VirtualHosts().Create(vh)
}
