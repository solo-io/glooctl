package upstream

import (
	"fmt"

	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func getCmd(opts *util.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name (optional)]",
		Short: "get upstream  or upstream list",
		Args:  cobra.MaximumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			output, _ := c.InheritedFlags().GetString("output")
			var name string
			if len(args) > 0 {
				name = args[0]
			}
			if err := runGet(sc, output, name); err != nil {
				fmt.Printf("Unable to get upstream %q\n", err)
				return
			}
		},
	}
	return cmd
}

func runGet(sc storage.Interface, output, name string) error {
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
		default:
			printSummaryList(upstreams)
		}
		return nil
	}

	upstream, err := sc.V1().Upstreams().Get(name)
	if err != nil {
		return err
	}
	switch output {
	case "json":
		printJSON(upstream)
	default:
		printYAML(upstream)
	}
	return nil
}
