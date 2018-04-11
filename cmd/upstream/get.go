package upstream

import (
	"fmt"
	"os"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
	"github.com/solo-io/glooctl/pkg/upstream"
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
		if len(upstreams) == 0 {
			fmt.Println("No upstreams found.")
			return nil
		}
		switch output {
		case "yaml":
			printYAMLList(upstreams)
		case "json":
			printJSONList(upstreams)
		case "template":
			return upstream.PrintTemplate(upstreams, tplt, os.Stdout)
		default:
			upstream.PrintTable(upstreams, os.Stdout)
		}
		return nil
	}

	u, err := sc.V1().Upstreams().Get(name)
	if err != nil {
		return err
	}
	switch output {
	case "json":
		printJSON(u)
	case "yaml":
		printYAML(u)
	case "template":
		upstream.PrintTemplate([]*v1.Upstream{u}, tplt, os.Stdout)
	default:
		upstream.PrintTable([]*v1.Upstream{u}, os.Stdout)
	}
	return nil
}
