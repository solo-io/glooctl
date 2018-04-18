package route

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	storage "github.com/solo-io/gloo/pkg/storage"
	proute "github.com/solo-io/glooctl/pkg/route"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func getCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get routes on a virtual host",
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			routes, err := runGet(sc, routeOpt.virtualhost, routeOpt.domain)
			if err != nil {
				fmt.Printf("Unable to get routes for %s: %q\n", routeOpt.domain, err)
				return
			}
			util.PrintList(routeOpt.output, "", routes,
				func(data interface{}, w io.Writer) error {
					proute.PrintTable(data.([]*v1.Route), w)
					return nil
				}, os.Stdout)
		},
	}
	return cmd
}

func runGet(sc storage.Interface, vhostname, domain string) ([]*v1.Route, error) {
	v, err := proute.VirtualHost(sc, vhostname, domain, false)
	if err != nil {
		return nil, err
	}
	fmt.Println("Using virtual host:", v.Name)
	return v.GetRoutes(), nil
}
