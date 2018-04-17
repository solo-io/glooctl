package vhost

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualhost"
	"github.com/spf13/cobra"
)

func updateCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			vh, err := runUpdate(sc, cliOpts.Filename)
			if err != nil {
				fmt.Printf("Unable to update virtual host %q\n", err)
				return
			}
			fmt.Println("Virtual host updated")
			util.Print(cliOpts.Output, "", vh, func(v interface{}, w io.Writer) error {
				virtualhost.PrintTable([]*v1.VirtualHost{v.(*v1.VirtualHost)}, w)
				return nil
			}, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&cliOpts.Filename, "filename", "f", "", "file to use to update virtual host")
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
	if vh.Metadata == nil {
		vh.Metadata = &v1.Metadata{}
	}
	vh.Metadata.ResourceVersion = existing.Metadata.GetResourceVersion()
	return sc.V1().VirtualHosts().Update(vh)
}
