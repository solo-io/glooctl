package vhost

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualhost"
	"github.com/spf13/cobra"
)

func updateCmd(opts *client.StorageOptions) *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := client.StorageClient(opts)
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
			util.Print(output, "", vh, func(v interface{}, w io.Writer) error {
				virtualhost.PrintTable([]*v1.VirtualHost{v.(*v1.VirtualHost)}, w)
				return nil
			}, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "file to use to update virtual host")
	cmd.MarkFlagFilename("filename", "yaml", "yml")
	cmd.MarkFlagRequired("filename")
	return cmd
}

func runUpdate(sc storage.Interface, filename string) (*v1.VirtualHost, error) {
	vh, err := parseFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load virtual host from %s", filename)
	}
	if err := defaultVHostValidation(vh); err != nil {
		return nil, err
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
