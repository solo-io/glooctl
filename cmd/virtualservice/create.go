package virtualservice

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
	"github.com/solo-io/glooctl/pkg/virtualservice"
	"github.com/spf13/cobra"
)

func createCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create virtual service",
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
				fmt.Printf("Unable to create virtual service %q\n", err)
				os.Exit(1)
			}
			util.Print(cliOpts.Output, "", vh, func(v interface{}, w io.Writer) error {
				virtualservice.PrintTable([]*v1.VirtualService{v.(*v1.VirtualService)}, w)
				return nil
			}, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&cliOpts.Filename, "filename", "f", "", "file to use to create virtual service")
	cmd.MarkFlagFilename("filename", "yaml", "yml")
	cmd.Flags().BoolVarP(&cliOpts.Interactive, "interactive", "i", true, "interactive mode")
	return cmd
}

func runCreate(sc storage.Interface, si dependencies.SecretStorage, opts *virtualservice.Options) (*v1.VirtualService, error) {
	var vs *v1.VirtualService
	if opts.Filename != "" {
		var err error
		vs, err = virtualservice.ParseFile(opts.Filename)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load virtual service from %s", opts.Filename)
		}
	} else {
		if !opts.Interactive {
			return nil, errors.New("no file specified and interactive mode turned off")
		}
		vs = &v1.VirtualService{}
		if err := virtualservice.Interactive(sc, si, vs); err != nil {
			return nil, err
		}
	}
	if err := virtualservice.DefaultVirtualServiceValidation(sc, vs); err != nil {
		return nil, err
	}
	return sc.V1().VirtualServices().Create(vs)
}
