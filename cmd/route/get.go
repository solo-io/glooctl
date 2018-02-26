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
			vhost, _ := c.InheritedFlags().GetString("vhost")

			routes, err := runGet(sc, vhost)
			if err != nil {
				fmt.Printf("Unable to get routes for %s: %q\n", vhost, err)
				return
			}
			output, _ := c.InheritedFlags().GetString("output")
			printRoutes(routes, output)
		},
	}
	return cmd
}

func runGet(sc storage.Interface, vhost string) ([]*v1.Route, error) {
	v, err := virtualHost(sc, vhost)
	if err != nil {
		return nil, err
	}
	fmt.Println("Using virtual host: ", vhost)
	return v.GetRoutes(), nil
}
