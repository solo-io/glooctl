package route

import (
	"fmt"
	"io"
	"os"

	"github.com/solo-io/glooctl/pkg/virtualservice"

	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/bootstrap"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/route"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func deleteCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a route",
		Long: `
There are three ways to select the route to delete:

  1. Using the interactive mode
  2. Passing the index of the route
  3. Specifying the route details of the matcher and destination
     in the CLI arguments or via YAML file.
	 While selecting routes to delete, glooctl matches routes based
	 on matcher and destintation only. It doesn't include extensions.`,
		Run: func(c *cobra.Command, args []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}

			runDelete(sc)
		},
	}
	cmd.Flags().BoolVarP(&routeOpt.Interactive, "interactive", "i", false, "interactive mode")
	cmd.Flags().IntVar(&routeOpt.Index, "index", 0, "index of the route you want to delete")
	return cmd
}

func runDelete(sc storage.Interface) {
	vs, err := virtualservice.VirtualService(sc, routeOpt.Virtualservice, routeOpt.Domain, false)
	if err != nil {
		fmt.Printf("Unable to get virtual service for routes: %q\n", err)
		os.Exit(1)
	}
	fmt.Println("Using virtual service:", vs.Name)
	routeOpt.Virtualservice = vs.Name
	routes := vs.GetRoutes()

	filtered, err := removeRoutes(sc, routes, routeOpt)
	if err != nil {
		fmt.Printf("Unable to delete routes: %q\n", err)
		os.Exit(1)
	}

	vs.Routes = filtered
	saved, err := save(sc, vs)
	if err != nil {
		fmt.Printf("Unable to save routes: %q\n", err)
		os.Exit(1)
	}
	util.PrintList(routeOpt.Output, "", saved,
		func(data interface{}, w io.Writer) error {
			route.PrintTable(data.([]*v1.Route), w)
			return nil
		}, os.Stdout)
}

func removeRoutes(sc storage.Interface, routes []*v1.Route, opts *route.RouteOption) ([]*v1.Route, error) {
	if opts.Interactive {
		result, err := route.SelectInteractive(routes, true)
		if err != nil {
			return nil, err
		}
		return result.NotSelected, nil
	}

	r, err := route.FromRouteOption(opts, sc)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get route")
	}
	return remove(routes, r)

}

func save(sc storage.Interface, virtualservice *v1.VirtualService) ([]*v1.Route, error) {
	updated, err := sc.V1().VirtualServices().Update(virtualservice)
	if err != nil {
		return nil, err
	}
	return updated.GetRoutes(), nil
}

func remove(routes []*v1.Route, route *v1.Route) ([]*v1.Route, error) {
	var updated []*v1.Route
	for _, r := range routes {
		if !match(route, r) {
			updated = append(updated, r)
		}
	}
	if len(routes) == len(updated) {
		return nil, errors.New("did not match any route")
	}
	return updated, nil
}

func match(left, right *v1.Route) bool {
	if !left.Matcher.Equal(right.Matcher) {
		return false
	}

	if left.GetSingleDestination() != nil && right.GetSingleDestination() != nil {
		return left.GetSingleDestination().Equal(right.GetSingleDestination())
	}

	// matching exact order of destintations
	if left.GetMultipleDestinations() != nil && right.GetMultipleDestinations() != nil {
		lm := left.GetMultipleDestinations()
		rm := right.GetMultipleDestinations()

		if len(lm) != len(rm) {
			return false
		}
		for i := range lm {
			if !lm[i].Equal(rm[i]) {
				return false
			}
		}
		return true
	}

	return false
}
