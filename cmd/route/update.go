package route

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"

	google_protobuf "github.com/gogo/protobuf/types"
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

			runUpdate(sc)
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

func runUpdate(sc storage.Interface) {
	v, err := proute.VirtualHost(sc, routeOpt.virtualhost, routeOpt.domain, false)
	if err != nil {
		fmt.Println("Unable to get virtual host for routes:", err)
		os.Exit(1)
	}
	fmt.Println("Using virtual host:", v.Name)
	routes := v.GetRoutes()
	updated, err := updateRoutes(sc, routes, routeOpt)
	if err != nil {
		fmt.Println("Unable to get updated route:", err)
		os.Exit(1)
	}

	v.Routes = updated
	if routeOpt.sort {
		sortRoutes(v.Routes)
	}

	saved, err := save(sc, v)
	if err != nil {
		fmt.Println("Unable to sav updated routes:", err)
		os.Exit(1)
	}
	util.PrintList(routeOpt.output, "", saved,
		func(data interface{}, w io.Writer) error {
			proute.PrintTable(data.([]*v1.Route), w)
			return nil
		}, os.Stdout)
}

func updateRoutes(sc storage.Interface, routes []*v1.Route, opts *routeOption) ([]*v1.Route, error) {
	if opts.interactive {
		selection, err := proute.SelectInteractive(routes, false)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get route")
		}
		if err := proute.RouteInteractive(sc, selection.Selected[0]); err != nil {
			return nil, err
		}
		return routes, nil // we have been working with pointers so it has changed the original route
	} else {
		route, err := route(opts, sc)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get route")
		}
		updated := make([]*v1.Route, len(routes))
		var matches []*v1.Route
		for i, r := range routes {
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
		return updated, nil
	}
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
