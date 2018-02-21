package vhost

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
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
			vh, err := runUpdate(sc, filename)
			if err != nil {
				fmt.Printf("Unable to update virtual host %q\n", err)
				return
			}
			fmt.Println("Virtual host updated")
			output, _ := c.InheritedFlags().GetString("output")
			if output == "yaml" {
				printYAML(vh)
			}
			if output == "json" {
				printJSON(vh)
			}
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "file to use to update virtual host")
	cmd.MarkFlagFilename("filename")
	cmd.MarkFlagRequired("filename")
	return cmd
}

func runUpdate(sc storage.Interface, filename string) (*v1.VirtualHost, error) {
	vh, err := parseFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load virtual host from %s", filename)
	}
	existing, err := sc.V1().VirtualHosts().Get(vh.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find existing virtual host %s", vh.Name)
	}
	// need to copy new into existing so that we get resource revision number
	copy(vh, existing)
	return sc.V1().VirtualHosts().Update(existing)
}

func copy(src *v1.VirtualHost, dst *v1.VirtualHost) {
	dst.Domains = src.Domains
	dst.Routes = src.Routes
	dst.SslConfig = src.SslConfig
}
