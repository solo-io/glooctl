package route

import (
	"fmt"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get routes on a virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			domain, _ := c.InheritedFlags().GetString("domain")

			routes, err := runGet(sc, domain)
			if err != nil {
				fmt.Printf("Unable to get routes for %s: %q\n", domain, err)
				return
			}
			output, _ := c.InheritedFlags().GetString("output")
			printRoutes(routes, output)
		},
	}
	return cmd
}

func runGet(sc storage.Interface, domain string) ([]*v1.Route, error) {
	v, created, err := virtualHost(sc, domain, false)
	if err != nil {
		return nil, err
	}
	if created {
		fmt.Println("Using newly created virtual host:", v.Name)
	} else {
		fmt.Println("Using virtual host:", v.Name)
	}
	return v.GetRoutes(), nil
}
