package vhost

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
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
				os.Exit(1)
			}
			si, err := client.SecretClient(opts)
			if err != nil {
				fmt.Printf("Unable to create secret client %q\n", err)
				os.Exit(1)
			}
			vh, err := runUpdate(sc, si, cliOpts)
			if err != nil {
				fmt.Printf("Unable to update virtual host %q\n", err)
				os.Exit(1)
			}
			util.Print(cliOpts.Output, "", vh, func(v interface{}, w io.Writer) error {
				virtualhost.PrintTable([]*v1.VirtualHost{v.(*v1.VirtualHost)}, w)
				return nil
			}, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&cliOpts.Filename, "filename", "f", "", "file to use to update virtual host")
	cmd.MarkFlagFilename("filename", "yaml", "yml")
	cmd.Flags().BoolVarP(&cliOpts.Interactive, "interactive", "i", true, "interactive mode")
	return cmd
}

func runUpdate(sc storage.Interface, si secret.SecretInterface, opts *virtualhost.Options) (*v1.VirtualHost, error) {
	var vh *v1.VirtualHost
	var existing *v1.VirtualHost
	if opts.Filename != "" {
		var err error
		vh, err := parseFile(opts.Filename)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load virtual host from %s", opts.Filename)
		}
		existing, err = sc.V1().VirtualHosts().Get(vh.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to find existing virtual host %s", vh.Name)
		}
	} else {
		var err error
		existing, err = virtualhost.SelectInteractive(sc)
		if err != nil {
			return nil, err
		}
		err = virtualhost.VirtualHostInteractive(sc, si, existing)
		if err != nil {
			return nil, err
		}
		vh = existing
	}
	if err := defaultVHostValidation(vh); err != nil {
		return nil, err
	}

	if vh.Metadata == nil {
		vh.Metadata = &v1.Metadata{}
	}
	vh.Metadata.ResourceVersion = existing.Metadata.GetResourceVersion()
	return sc.V1().VirtualHosts().Update(vh)
}
