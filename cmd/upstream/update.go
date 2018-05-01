package upstream

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	"github.com/solo-io/gloo/pkg/bootstrap/secretstorage"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/gloo/pkg/storage/dependencies"
	"github.com/solo-io/glooctl/pkg/upstream"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func updateCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update upstreams",
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
			u, err := runUpdate(sc, si, cliOpts)
			if err != nil {
				fmt.Printf("Unable to update upstream %q\n", err)
				os.Exit(1)
			}
			fmt.Println("Upstream updated")
			util.Print(cliOpts.Output, "", u,
				func(data interface{}, w io.Writer) error {
					upstream.PrintTable([]*v1.Upstream{data.(*v1.Upstream)}, w)
					return nil
				}, os.Stdout)
		},
	}
	cmd.Flags().StringVarP(&cliOpts.Filename, "filename", "f", "", "file to use to update upstream")
	cmd.MarkFlagFilename("filename", "yaml", "yml")

	cmd.Flags().BoolVarP(&cliOpts.Interactive, "interactive", "i", true, "interactive mode")
	return cmd
}

func runUpdate(sc storage.Interface, si dependencies.SecretStorage, opts *upstream.Options) (*v1.Upstream, error) {
	var u *v1.Upstream
	var existing *v1.Upstream
	var err error
	if opts.Filename != "" {
		u, err = upstream.ParseFile(opts.Filename)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load Upstream from %s", opts.Filename)
		}
		valid, message := upstream.Validate(sc, si, u)
		if !valid {
			return nil, fmt.Errorf("invalid upstream: %s", message)
		}
		if message != "" {
			fmt.Println("Warning:", message)
		}
		existing, err = sc.V1().Upstreams().Get(u.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to find existing Upstream %s", u.Name)
		}
	} else {
		if !opts.Interactive {
			return nil, errors.New("no file specified and interactive mode turned off")
		}

		existing, err = upstream.SelectInteractive(sc)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get upstream to edit")
		}
		err = upstream.Interactive(sc, si, existing)
		if err != nil {
			return nil, err
		}
		u = existing
	}
	if u.Metadata == nil {
		u.Metadata = &v1.Metadata{}
	}
	u.Metadata.ResourceVersion = existing.Metadata.GetResourceVersion()
	return sc.V1().Upstreams().Update(u)
}
