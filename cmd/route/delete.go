package route

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
	proute "github.com/solo-io/glooctl/pkg/route"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func deleteCmd(opts *client.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a route",
		Long: `
Delete a route based on either the definition in the YAML file
or based on the route matcher and destintation provided in the CLI.

While selecting routes to delete, glooctl matches routes based on 
matcher and destintation only. It doesn't include extensions.`,
		Run: func(c *cobra.Command, args []string) {
			sc, err := client.StorageClient(opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}

			runDelete(sc)
		},
	}
	cmd.Flags().BoolVarP(&routeOpt.interactive, "interactive", "i", false, "interactive mode")
	return cmd
}

func runDelete(sc storage.Interface) {
	vh, err := proute.VirtualHost(sc, routeOpt.virtualhost, routeOpt.domain, false)
	if err != nil {
		fmt.Printf("Unable to get virtualhost for routes: %q\n", err)
		os.Exit(1)
	}
	fmt.Println("Using virtual host:", vh.Name)
	routes := vh.GetRoutes()

	filtered, err := removeRoutes(sc, routes, routeOpt)
	if err != nil {
		fmt.Printf("Unable to delete routes: %q\n", err)
		os.Exit(1)
	}

	vh.Routes = filtered
	saved, err := save(sc, vh)
	if err != nil {
		fmt.Printf("Unable to save routes: %q\n", err)
		os.Exit(1)
	}
	util.PrintList(routeOpt.output, "", saved,
		func(data interface{}, w io.Writer) error {
			proute.PrintTable(data.([]*v1.Route), w)
			return nil
		}, os.Stdout)
}

func removeRoutes(sc storage.Interface, routes []*v1.Route, opts *routeOption) ([]*v1.Route, error) {
	if opts.interactive {
		result, err := proute.SelectInteractive(routes, true)
		if err != nil {
			return nil, err
		}
		return result.NotSelected, nil
	} else {
		route, err := route(opts, sc)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get route")
		}
		return remove(routes, route)
	}
}

func save(sc storage.Interface, virtualhost *v1.VirtualHost) ([]*v1.Route, error) {
	updated, err := sc.V1().VirtualHosts().Update(virtualhost)
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
