package upstream

// FIXME - replace kube secret interface with secrets client
import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	"github.com/solo-io/gloo/pkg/bootstrap/secretstorage"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/gloo/pkg/storage/dependencies"
	"github.com/solo-io/glooctl/pkg/upstream"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func createCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create upstream",
		Long:  "Create an upstream based on the file or interactively if no file is specified",
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
			u, err := runCreate(sc, si, cliOpts)
			if err != nil {
				fmt.Printf("Unable to create upstream %q\n", err)
				os.Exit(1)
			}
			util.Print(cliOpts.Output, cliOpts.Template, u,
				func(data interface{}, w io.Writer) error {
					upstream.PrintTable([]*v1.Upstream{data.(*v1.Upstream)}, w)
					return nil
				}, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&cliOpts.Filename, "filename", "f", "", "file to use to create upstream")
	cmd.MarkFlagFilename("filename", "yaml", "yml")

	cmd.Flags().BoolVarP(&cliOpts.Interactive, "interactive", "i", true, "interacitve mode")
	return cmd
}

func runCreate(sc storage.Interface, si dependencies.SecretStorage, opts *upstream.Options) (*v1.Upstream, error) {
	var u *v1.Upstream
	if opts.Filename != "" {
		var err error
		u, err = upstream.ParseFile(opts.Filename)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load Upstream from %s", opts.Filename)
		}

		valid, message := upstream.Validate(sc, si, u)
		if !valid {
			return nil, fmt.Errorf("invalid upstream: %s", message)
		}
		// add verbose mode to disable this normally
		if message != "" {
			fmt.Println("Warning:", message)
		}
	} else {
		// if not file then interactive unless explicitly not set to true
		if !opts.Interactive {
			return nil, errors.New("no file specified and interactive mode turned off")
		}
		u = &v1.Upstream{}
		err := upstream.Interactive(sc, si, u)
		if err != nil {
			return nil, err
		}
	}
	return sc.V1().Upstreams().Create(u)
}
