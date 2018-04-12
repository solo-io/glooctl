package vhost

import (
	"fmt"
	"os"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
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
		if len(v) == 0 {
			fmt.Println("No virtual hosts found.")
			return nil
		}
		switch output {
		case "yaml":
			printYAMLList(v)
		case "json":
			printJSONList(v)
		case "template":
			return virtualhost.PrintTemplate(v, tplt, os.Stdout)
		default:
			virtualhost.PrintTable(v, os.Stdout)
		}
		return nil
	}

	v, err := sc.V1().VirtualHosts().Get(name)
	if err != nil {
		return err
	}
	switch output {
	case "json":
		printJSON(v)
	case "yaml":
		printYAML(v)
	case "template":
		return virtualhost.PrintTemplate([]*v1.VirtualHost{v}, tplt, os.Stdout)
	default:
		virtualhost.PrintTable([]*v1.VirtualHost{v}, os.Stdout)
	}
	return nil
}
