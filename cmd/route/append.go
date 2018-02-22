package route

import (
	"fmt"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func appendCmd() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "append",
		Short: "append a route to on a virtual host's routes",
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			vhost, _ := c.InheritedFlags().GetString("vhost")
			routes, err := runUpdate(sc, vhost, filename)
			if err != nil {
				fmt.Printf("Unable to get routes for %s: %q\n", vhost, err)
				return
			}
			output, _ := c.InheritedFlags().GetString("output")
			printRoutes(routes, output)
		},
	}
	cmd.Flags().StringVarP(&filename, "filename", "f", "", "file to route to append")
	cmd.MarkFlagFilename("filename")
	cmd.MarkFlagRequired("filename")
	return cmd
}

func runUpdate(sc storage.Interface, vhost, filename string) ([]*v1.Route, error) {
	route, err := parseFile(filename)
	if err != nil {
		return nil, err
	}
	v, err := sc.V1().VirtualHosts().Get(vhost)
	if err != nil {
		return nil, err
	}
	v.Routes = append(v.GetRoutes(), route)
	updated, err := sc.V1().VirtualHosts().Update(v)
	if err != nil {
		return nil, err
	}
	return updated.GetRoutes(), nil
}
