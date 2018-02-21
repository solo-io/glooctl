package vhost

import (
	"fmt"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name (optional)]",
		Short: "get virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
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
				fmt.Printf("Unable to get virtual host %q\n", err)
				return
			}
		},
	}
	return cmd
}

func runGet(sc storage.Interface, output, name string) error {
	if name == "" {
		v, err := sc.V1().VirtualHosts().List()
		if err != nil {
			return err
		}
		if len(v) == 0 {
			fmt.Println("No virtual hosts found.")
			return nil
		}
		switch output {
		case "yaml":
			printYAMLList(v)
		case "json":
			printJSONList(v)
		default:
			printSummaryList(v)
		}
	}

	v, err := sc.V1().VirtualHosts().Get(name)
	if err != nil {
		return err
	}
	switch output {
	case "json":
		printJSON(v)
	default:
		printYAML(v)
	}
	return nil
}
