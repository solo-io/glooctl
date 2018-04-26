package vhost

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/solo-io/gloo/pkg/bootstrap/secretstorage"

	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualhost"
	"github.com/spf13/cobra"
)

func createCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				os.Exit(1)
			}
			si, err := secretstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create secret client %q\n", err)
				os.Exit(1)
			}
			vh, err := runCreate(sc, si, cliOpts)
			if err != nil {
				fmt.Printf("Unable to create virtual host %q\n", err)
				os.Exit(1)
			}
			util.Print(cliOpts.Output, "", vh, func(v interface{}, w io.Writer) error {
				virtualhost.PrintTable([]*v1.VirtualHost{v.(*v1.VirtualHost)}, w)
				return nil
			}, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&cliOpts.Filename, "filename", "f", "", "file to use to create virtual host")
	cmd.MarkFlagFilename("filename", "yaml", "yml")
	cmd.Flags().BoolVarP(&cliOpts.Interactive, "interactive", "i", true, "interactive mode")
	return cmd
}

func runCreate(sc storage.Interface, si dependencies.SecretStorage, opts *virtualhost.Options) (*v1.VirtualHost, error) {
	var vh *v1.VirtualHost
	if opts.Filename != "" {
		var err error
		vh, err = parseFile(opts.Filename)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load virtual host from %s", opts.Filename)
		}
	} else {
		if !opts.Interactive {
			return nil, errors.New("no file specified and interactive mode turned off")
		}
		vh = &v1.VirtualHost{}
		if err := virtualhost.Interactive(sc, si, vh); err != nil {
			return nil, err
		}
	}
	if err := defaultVHostValidation(vh); err != nil {
		return nil, err
	}
	return sc.V1().VirtualHosts().Create(vh)
}
