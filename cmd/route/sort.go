package route

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	storage "github.com/solo-io/gloo/pkg/storage"
	proute "github.com/solo-io/glooctl/pkg/route"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualservice"
	"github.com/spf13/cobra"
)

func sortCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sort",
		Short: "sort routes to have the longest route first",
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				os.Exit(1)
			}

			routes, err := runSort(sc, routeOpt.virtualservice, routeOpt.domain)
			if err != nil {
				fmt.Println("Unable to sort routes", err)
				os.Exit(1)
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

func runSort(sc storage.Interface, vservicename, domain string) ([]*v1.Route, error) {
	v, err := virtualservice.VirtualService(sc, vservicename, domain, false)
	if err != nil {
		return nil, err
	}
	fmt.Println("Using virtual service:", v.Name)
	proute.SortRoutes(v.Routes)
	updated, err := sc.V1().VirtualServices().Update(v)
	if err != nil {
		return nil, err
	}
	return updated.GetRoutes(), nil
}
