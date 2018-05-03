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

func createCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a route to a destination",
		Long: `
Create a route. The route, with its matcher and destination, can be provided
using a file or by specifying one of the matcher and a destintation using
the flags.`,
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			var r *v1.Route
			if routeOpt.interactive {
				r = &v1.Route{}
				err = proute.Interactive(sc, r)
			} else {
				r, err = route(routeOpt, sc)
			}
			if err != nil {
				fmt.Printf("Unable to get route %q\n", err)
				return
			}

			routes, err := runCreate(sc, routeOpt.virtualservice, routeOpt.domain, r, routeOpt.sort)
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
	kube := routeOpt.route.kube
	flags := cmd.Flags()
	flags.StringVar(&kube.name, flagKubeName, "", "kubernetes service name")
	flags.StringVar(&kube.namespace, flagKubeNamespace, "", "kubernetes service namespace")
	flags.IntVar(&kube.port, flagKubePort, 0, "kubernetes service port")
	flags.BoolVar(&routeOpt.sort, "sort", false, "sort the routes after appending the new route")
	flags.BoolVarP(&routeOpt.interactive, "interactive", "i", false, "interactive mode")
	return cmd
}

func runCreate(sc storage.Interface, vservicename, domain string, route *v1.Route, sort bool) ([]*v1.Route, error) {
	v, err := virtualservice.VirtualService(sc, vservicename, domain, true)
	if err != nil {
		return nil, err
	}
	fmt.Println("Using virtual service:", v.Name)
	v.Routes = append(v.GetRoutes(), route)
	if sort {
		sortRoutes(v.Routes)
	}
	updated, err := sc.V1().VirtualServices().Update(v)
	if err != nil {
		return nil, err
	}
	return updated.GetRoutes(), nil
}
