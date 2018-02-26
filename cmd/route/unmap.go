package route

import (
	"fmt"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func unmapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unmap",
		Short: "unmap a route",
		Long: `
Unmap a route based on either the definition in the YAML file
or based on the route matcher and destintation provided in the CLI.

glooctl matches routes based on matcher and destintation only. It
doesn't include extensions.`,
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}

			flags := c.InheritedFlags()
			vhost, _ := flags.GetString("vhost")
			route, err := route(flags)
			if err != nil {
				fmt.Printf("Unable to get route %q\n", err)
				return
			}
			routes, err := runUnmap(sc, vhost, route)
			if err != nil {
				fmt.Printf("Unable to get route for %s: %q\n", vhost, err)
				return
			}
			output, _ := flags.GetString("output")
			printRoutes(routes, output)
		},
	}
	return cmd
}

func runUnmap(sc storage.Interface, vhost string, route *v1.Route) ([]*v1.Route, error) {
	v, err := virtualHost(sc, vhost)
	if err != nil {
		return nil, err
	}
	fmt.Println("Using virtual host: ", vhost)
	v.Routes = remove(v.GetRoutes(), route)
	updated, err := sc.V1().VirtualHosts().Update(v)
	if err != nil {
		return nil, err
	}
	return updated.GetRoutes(), nil
}

func remove(routes []*v1.Route, route *v1.Route) []*v1.Route {
	var updated []*v1.Route
	for _, r := range routes {
		if !match(route, r) {
			fmt.Println("didint' match")
			updated = append(updated, r)
		}
	}
	return updated
}

func match(left, right *v1.Route) bool {
	fmt.Println("right: ")
	printYAML(right)
	if !left.Matcher.Equal(right.Matcher) {
		fmt.Println("matcher didn't match")
		return false
	}

	if left.GetSingleDestination() != nil && right.GetSingleDestination() != nil {
		fmt.Println("comparing single destination")
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
