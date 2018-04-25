package vhost

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/editor"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualhost"
	"github.com/spf13/cobra"
)

func editCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [name]",
		Short: "edit a virtual host",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				os.Exit(1)
			}

			v, err := runEdit(sc, args[0])
			if err != nil {
				fmt.Printf("Unable to edit virtual host %s: %q\n", args[0], err)
				os.Exit(1)
			}
			fmt.Printf("Virtual host %s updated\n", args[0])

			util.Print(cliOpts.Output, "", v, func(v interface{}, w io.Writer) error {
				virtualhost.PrintTable([]*v1.VirtualHost{v.(*v1.VirtualHost)}, w)
				return nil
			}, os.Stdout)
		},
	}
	return cmd
}

func runEdit(sc storage.Interface, name string) (*v1.VirtualHost, error) {
	if name == "" {
		return nil, fmt.Errorf("missing name of virtual host to edit")
	}

	v, err := sc.V1().VirtualHosts().Get(name)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get virtual host "+name)
	}
	edit := editor.NewDefaultEditor([]string{"EDITOR", "GLOO_EDITOR"})

	f, err := ioutil.TempFile("", "thetool-virtualhost")
	if err != nil {
		return nil, errors.Wrap(err, "unable to create temporary file")
	}
	err = util.PrintYAML(v, f)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write out virtualhost for editting")
	}
	defer func() {
		if errRemove := os.Remove(f.Name()); errRemove != nil {
			fmt.Fprintln(os.Stderr, "unable to remove temporary file", f.Name())
		}
	}()

	err = edit.Launch(f.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to edit virtualhost")
	}
	updated, err := parseFile(f.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse virtual host "+name)
	}
	return sc.V1().VirtualHosts().Update(updated)
}
