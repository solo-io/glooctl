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
	"github.com/solo-io/glooctl/pkg/route"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualservice"
	"github.com/spf13/cobra"
)

var (
	// represents the new route defintion for update
	oldRouteOpt = &route.RouteOption{Route: &route.RouteDetail{Kube: &route.KubeUpstream{}}}
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
	kube := routeOpt.Route.Kube
	flags := cmd.Flags()
	flags.StringVar(&kube.Name, flagKubeName, "", "kubernetes service name")
	flags.StringVar(&kube.Namespace, flagKubeNamespace, "", "kubernetes service namespace")
	flags.IntVar(&kube.Port, flagKubePort, 0, "kubernetes service port")
	flags.BoolVar(&routeOpt.Sort, "sort", false, "sort the routes after appending the new route")
	flags.BoolVarP(&routeOpt.Interactive, "interactive", "i", false, "interactive mode")

	setupOldRouteParams(cmd)
	return cmd
}

func setupOldRouteParams(cmd *cobra.Command) {
	flags := cmd.Flags()
	r := oldRouteOpt.Route
	flags.StringVar(&r.Event, "old-"+flagEvent, "", "event type to match")
	flags.StringVar(&r.PathExact, "old-"+flagPathExact, "", "exact path to match")
	flags.StringVar(&r.PathRegex, "old-"+flagPathRegex, "", "path regex to match")
	flags.StringVar(&r.PathPrefix, "old-"+flagPathPrefix, "", "path prefix to match")
	flags.StringVar(&r.Verb, "old-"+flagMethod, "", "HTTP method to match")
	flags.StringVar(&r.Headers, "old-"+flagHeaders, "", "header to match")
	flags.StringVar(&r.Upstream, "old-"+flagUpstream, "", "desitnation upstream")
	flags.StringVar(&r.Function, "old-"+flagFunction, "", "destination function")

	kube := r.Kube
	flags.StringVar(&kube.Name, "old-"+flagKubeName, "", "kubernetes service name")
	flags.StringVar(&kube.Namespace, "old-"+flagKubeNamespace, "", "kubernetes service namespace")
	flags.IntVar(&kube.Port, "old-"+flagKubePort, 0, "kubernetes service port")

	// auto complete
	annotate(cmd.Flag("old-"+flagMethod), "__glooctl_route_http_methods")
	annotate(cmd.Flag("old-"+flagUpstream), "__glooctl_get_upstreams")
	annotate(cmd.Flag("old-"+flagFunction), "__glooctl_get_functions")
}

func runUpdate(sc storage.Interface) {
	v, err := virtualservice.VirtualService(sc, routeOpt.Virtualservice, routeOpt.Domain, false)
	if err != nil {
		fmt.Println("Unable to get virtual service for routes:", err)
		os.Exit(1)
	}
	fmt.Println("Using virtual service:", v.Name)
	routes := v.GetRoutes()
	updated, err := updateRoutes(sc, routes, routeOpt, oldRouteOpt)
	if err != nil {
		fmt.Println("Unable to get updated route:", err)
		os.Exit(1)
	}

	v.Routes = updated
	if routeOpt.Sort {
		route.SortRoutes(v.Routes)
	}

	saved, err := save(sc, v)
	if err != nil {
		fmt.Println("Unable to sav updated routes:", err)
		os.Exit(1)
	}
	util.PrintList(routeOpt.Output, "", saved,
		func(data interface{}, w io.Writer) error {
			route.PrintTable(data.([]*v1.Route), w)
			return nil
		}, os.Stdout)
}

func updateRoutes(sc storage.Interface, routes []*v1.Route, opts, oldOpts *route.RouteOption) ([]*v1.Route, error) {
	if opts.Interactive {
		selection, err := route.SelectInteractive(routes, false)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get route")
		}
		if err := route.Interactive(sc, selection.Selected[0]); err != nil {
			return nil, err
		}
		return routes, nil // we have been working with pointers so it has changed the original route
	}

	newRoute, err := route.FromRouteOption(opts, sc)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get new route")
	}
	oldRoute, err := route.FromRouteOption(oldOpts, sc)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get old route")
	}
	updated := make([]*v1.Route, len(routes))
	var matches []*v1.Route
	for i, r := range routes {
		if match(oldRoute, r) {
			matches = append(matches, r)
			newRoute.Extensions = mergeExtensions(newRoute, r)
			updated[i] = newRoute
			continue
		}
		updated[i] = r
	}
	if len(matches) == 0 {
		return nil, errors.New("could not find a route for the specified matcher and destination")
	}
	if len(matches) > 1 {
		return nil, errors.New("found more than one route for the specified matcher and destination")
	}
	return updated, nil
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
