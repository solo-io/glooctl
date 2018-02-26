package route

import (
	"fmt"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a route to a destination",
		Long: `
Create a route. The route, with its matcher and destination, can be provided
using a file or by specifying one of the matcher and a destintation using
the flags.`,
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			flags := c.InheritedFlags()
			domain, _ := flags.GetString("domain")
			route, err := route(flags)
			if err != nil {
				fmt.Printf("Unable to get route %q\n", err)
				return
			}
			routes, err := runCreate(sc, domain, route)
			if err != nil {
				fmt.Printf("Unable to get routes for %s: %q\n", domain, err)
				return
			}
			output, _ := flags.GetString("output")
			printRoutes(routes, output)
		},
	}
	return cmd
}

func runCreate(sc storage.Interface, domain string, route *v1.Route) ([]*v1.Route, error) {
	v, created, err := virtualHost(sc, domain, true)
	if err != nil {
		return nil, err
	}
	if created {
		fmt.Println("Using newly virtual host: ", v.Name)
	} else {
		fmt.Println("Using virtual host: ", v.Name)
	}
	v.Routes = append(v.GetRoutes(), route)
	updated, err := sc.V1().VirtualHosts().Update(v)
	if err != nil {
		return nil, err
	}
	return updated.GetRoutes(), nil
}
