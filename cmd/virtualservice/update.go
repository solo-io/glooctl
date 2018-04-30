package virtualservice

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/solo-io/gloo/pkg/bootstrap/secretstorage"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualservice"
	"github.com/spf13/cobra"
)

func updateCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update virtual service",
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
			vh, err := runUpdate(sc, si, cliOpts)
			if err != nil {
				fmt.Printf("Unable to update virtual service %q\n", err)
				os.Exit(1)
			}
			util.Print(cliOpts.Output, "", vh, func(v interface{}, w io.Writer) error {
				virtualservice.PrintTable([]*v1.VirtualService{v.(*v1.VirtualService)}, w)
				return nil
			}, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&cliOpts.Filename, "filename", "f", "", "file to use to update virtual service")
	cmd.MarkFlagFilename("filename", "yaml", "yml")
	cmd.Flags().BoolVarP(&cliOpts.Interactive, "interactive", "i", true, "interactive mode")
	return cmd
}

func runUpdate(sc storage.Interface, si dependencies.SecretStorage, opts *virtualservice.Options) (*v1.VirtualService, error) {
	var vh *v1.VirtualService
	var existing *v1.VirtualService
	if opts.Filename != "" {
		var err error
		vh, err = parseFile(opts.Filename)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load virtual service from %s", opts.Filename)
		}
		existing, err = sc.V1().VirtualServices().Get(vh.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to find existing virtual service %s", vh.Name)
		}
	} else {
		var err error
		existing, err = virtualservice.SelectInteractive(sc)
		if err != nil {
			return nil, err
		}
		err = virtualservice.Interactive(sc, si, existing)
		if err != nil {
			return nil, err
		}
		vh = existing
	}
	if err := defaultVirtualServiceValidation(vh); err != nil {
		return nil, err
	}

	if vh.Metadata == nil {
		vh.Metadata = &v1.Metadata{}
	}
	vh.Metadata.ResourceVersion = existing.Metadata.GetResourceVersion()
	return sc.V1().VirtualServices().Update(vh)
}
