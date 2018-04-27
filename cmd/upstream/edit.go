package upstream

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/editor"
	"github.com/solo-io/glooctl/pkg/upstream"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func editCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [name]",
		Short: "edit an upstream",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}

			u, err := runEdit(sc, args[0])
			if err != nil {
				fmt.Printf("Unable to edit upstream %s: %q\n", args[0], err)
				return
			}
			fmt.Printf("Upstream %s updated\n", args[0])
			util.Print(cliOpts.Output, cliOpts.Template, u,
				func(data interface{}, w io.Writer) error {
					upstream.PrintTable([]*v1.Upstream{data.(*v1.Upstream)}, w)
					return nil
				}, os.Stdout)
		},
	}
	return cmd
}

func runEdit(sc storage.Interface, name string) (*v1.Upstream, error) {
	if name == "" {
		return nil, fmt.Errorf("missing name of upstream to edit")
	}

	u, err := sc.V1().Upstreams().Get(name)
	edit := editor.NewDefaultEditor([]string{"EDITOR", "GLOO_EDITOR"})
	if err != nil {
		return nil, errors.Wrap(err, "unable to get upstream "+name)
	}

	f, err := ioutil.TempFile("", "thetool-upstream")
	if err != nil {
		return nil, errors.Wrap(err, "unable to create temporary file")
	}
	err = util.PrintYAML(u, f)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write out upstream for editting")
	}
	defer os.Remove(f.Name())
	err = edit.Launch(f.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to edit upstream")
	}
	updated, err := upstream.ParseFile(f.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse upstream")
	}
	return sc.V1().Upstreams().Update(updated)
}
