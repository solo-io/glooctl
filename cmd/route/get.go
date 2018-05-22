package route

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/glooctl/pkg/virtualservice"

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
		Short: "get routes on a virtual service",
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			routes, err := runGet(sc, routeOpt.Virtualservice, routeOpt.Domain)
			if err != nil {
				fmt.Printf("Unable to get routes for %s: %q\n", routeOpt.Domain, err)
				return
			}
			util.PrintList(routeOpt.Output, "", routes,
				func(data interface{}, w io.Writer) error {
					proute.PrintTable(data.([]*v1.Route), w)
					return nil
				}, os.Stdout)
		},
	}
	return cmd
}

func runGet(sc storage.Interface, vservicename, domain string) ([]*v1.Route, error) {
	v, err := virtualservice.VirtualService(sc, vservicename, domain, false)
	if err != nil {
		return nil, err
	}
	fmt.Println("Using virtual service:", v.Name)
	return v.GetRoutes(), nil
}
