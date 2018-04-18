package route

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"

	google_protobuf "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	storage "github.com/solo-io/gloo/pkg/storage"
	proute "github.com/solo-io/glooctl/pkg/route"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func updateCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update a route",
		Long: `
Update a route based on either the definition in the YAML file
or based on the route matcher and destination provided in the CLI.

While selecting route to update, glooctl matches routes based on
matcher and destination only. It doesn't include extensions.`,
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}

			route, err := route(routeOpt, sc)
			if err != nil {
				fmt.Printf("Unable to get route %q\n", err)
				return
			}
			routes, err := runUpdate(sc, routeOpt.virtualhost, routeOpt.domain, route, routeOpt.sort)
			if err != nil {
				fmt.Printf("Unable to get route for %s: %q\n", routeOpt.virtualhost, err)
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
	return cmd
}

func runUpdate(sc storage.Interface, vhostname, domain string, route *v1.Route, sort bool) ([]*v1.Route, error) {
	v, err := proute.VirtualHost(sc, vhostname, domain, false)
	if err != nil {
		return nil, err
	}
	fmt.Println("Using virtual host:", v.Name)

	existing := v.GetRoutes()
	updated := make([]*v1.Route, len(existing))
	var matches []*v1.Route
	for i, r := range existing {
		if match(route, r) {
			matches = append(matches, r)
			route.Extensions = mergeExtensions(route, r)
			updated[i] = route
			continue
		}
		updated[i] = r
	}
	if len(matches) == 0 {
		return nil, errors.New("could not find a route for the specified matcher and destination.")
	}
	if len(matches) > 1 {
		return nil, errors.New("found more than one route for the specified matcher and destination")
	}

	v.Routes = updated
	if sort {
		sortRoutes(v.Routes)
	}

	vh, err := sc.V1().VirtualHosts().Update(v)
	if err != nil {
		return nil, err
	}
	return vh.GetRoutes(), nil
}

func mergeExtensions(route, old *v1.Route) *google_protobuf.Struct {
	if old.Extensions == nil || old.Extensions.Fields == nil {
		return route.Extensions
	}

	if route.Extensions == nil || route.Extensions.Fields == nil {
		return old.Extensions
	}

	for k, v := range route.Extensions.Fields {
		old.Extensions.Fields[k] = v
	}

	return old.Extensions
}
