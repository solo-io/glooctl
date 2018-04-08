package upstream

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/editor"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func editCmd(opts *util.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [name]",
		Short: "edit an upstream",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(opts)
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

			output, _ := c.InheritedFlags().GetString("output")
			if output == "yaml" {
				printYAML(u)
			}
			if output == "json" {
				printJSON(u)
			}
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
	if err := writeYAML(u, f); err != nil {
		return nil, errors.Wrap(err, "unable to write out upstream for editting")
	}
	defer os.Remove(f.Name())
	if err := edit.Launch(f.Name()); err != nil {
		return nil, errors.Wrap(err, "unable to edit upstream")
	}
	updated, err := parseFile(f.Name())
	return sc.V1().Upstreams().Update(updated)
}
