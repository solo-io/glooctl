package vhost

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualhost"
	"github.com/spf13/cobra"
)

func getCmd(opts *client.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name (optional)]",
		Short: "get virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := client.StorageClient(opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			if output == "template" && tplt == "" {
				fmt.Println("Must provide template when setting output as template")
				return
			}
			var name string
			if len(args) > 0 {
				name = args[0]
			}
			if err := runGet(sc, output, tplt, name); err != nil {
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
