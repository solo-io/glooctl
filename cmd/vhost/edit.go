package vhost

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/editor"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func editCmd(opts *util.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [name]",
		Short: "edit a virtual host",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}

			v, err := runEdit(sc, args[0])
			if err != nil {
				fmt.Printf("Unable to edit virtual host %s: %q\n", args[0], err)
				return
			}
			fmt.Printf("Virtual host %s updated\n", args[0])

			output, _ := c.InheritedFlags().GetString("output")
			if output == "yaml" {
				printYAML(v)
			}
			if output == "json" {
				printJSON(v)
			}
		},
	}
	return cmd
}

func runEdit(sc storage.Interface, name string) (*v1.VirtualHost, error) {
	if name == "" {
		return nil, fmt.Errorf("missing name of virtual host to edit")
	}

	v, err := sc.V1().VirtualHosts().Get(name)
	edit := editor.NewDefaultEditor([]string{"EDITOR", "GLOO_EDITOR"})
	if err != nil {
		return nil, errors.Wrap(err, "unable to get virtual host "+name)
	}

	f, err := ioutil.TempFile("", "thetool-virtualhost")
	if err != nil {
		return nil, errors.Wrap(err, "unable to create temporary file")
	}
	if err := writeYAML(v, f); err != nil {
		return nil, errors.Wrap(err, "unable to write out virtualhost for editting")
	}
	defer os.Remove(f.Name())
	if err := edit.Launch(f.Name()); err != nil {
		return nil, errors.Wrap(err, "unable to edit virtualhost")
	}
	updated, err := parseFile(f.Name())
	return sc.V1().VirtualHosts().Update(updated)
}
