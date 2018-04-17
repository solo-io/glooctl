package vhost

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualhost"
	"github.com/spf13/cobra"
)

func getCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name (optional)]",
		Short: "get virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			if cliOpts.Output == "template" && cliOpts.Template == "" {
				fmt.Println("Must provide template when setting output as template")
				return
			}
			var name string
			if len(args) > 0 {
				name = args[0]
			}
			if err := runGet(sc, cliOpts.Output, cliOpts.Template, name); err != nil {
				fmt.Printf("Unable to get virtual host %q\n", err)
				return
			}
		},
	}
	return cmd
}

func runGet(sc storage.Interface, output, tplt, name string) error {
	if name == "" {
		v, err := sc.V1().VirtualHosts().List()
		if err != nil {
			return err
		}
		return util.PrintList(output, tplt, v,
			func(data interface{}, w io.Writer) error {
				virtualhost.PrintTable(data.([]*v1.VirtualHost), w)
				return nil
			}, os.Stdout)
	}

	v, err := sc.V1().VirtualHosts().Get(name)
	if err != nil {
		return err
	}
	return util.Print(output, tplt, v, func(v interface{}, w io.Writer) error {
		virtualhost.PrintTable([]*v1.VirtualHost{v.(*v1.VirtualHost)}, w)
		return nil
	}, os.Stdout)
}
