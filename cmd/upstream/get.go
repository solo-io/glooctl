package upstream

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
	"github.com/solo-io/glooctl/pkg/upstream"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func getCmd(opts *client.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name (optional)]",
		Short: "get upstream  or upstream list",
		Args:  cobra.MaximumNArgs(1),
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
				fmt.Printf("Unable to get upstream %q\n", err)
				return
			}
		},
	}
	return cmd
}

func runGet(sc storage.Interface, output, tplt, name string) error {
	if name == "" {
		upstreams, err := sc.V1().Upstreams().List()
		if err != nil {
			return err
		}
		return util.PrintList(output, tplt, upstreams,
			func(data interface{}, w io.Writer) error {
				upstream.PrintTable(data.([]*v1.Upstream), w)
				return nil
			}, os.Stdout)
	}

	u, err := sc.V1().Upstreams().Get(name)
	if err != nil {
		return err
	}
	return util.Print(output, tplt, u,
		func(data interface{}, w io.Writer) error {
			upstream.PrintTable([]*v1.Upstream{data.(*v1.Upstream)}, w)
			return nil
		}, os.Stdout)
}
