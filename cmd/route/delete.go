package route

import (
	"fmt"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func deleteCmd() *cobra.Command {
	var filename string
	var event string
	var pathExact string
	var pathPrefix string
	var pathRegex string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a route from a virtual host's routes",
		Long: `
Deletes a route based on either the definition in the YAML file
or based on the router matcher provided in the CLI. The route is
created using one of the flags event, path-exact, path-regex, 
path-prefix or filename`,
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}

			route, err := route(event, pathExact, pathRegex, pathPrefix, filename)
			if err != nil {
				fmt.Printf("Unable to get route %q\n", err)
				return
			}
			vhost, _ := c.InheritedFlags().GetString("vhost")
			routes, err := runDelete(sc, vhost, route)
			if err != nil {
				fmt.Printf("Unable to get routes for %s: %q\n", vhost, err)
				return
			}
			output, _ := c.InheritedFlags().GetString("output")
			printRoutes(routes, output)
		},
	}
	cmd.Flags().StringVarP(&filename, "filename", "f", "", "file to route to delete")
	cmd.Flags().StringVarP(&event, "event", "e", "", "event matcher")
	cmd.Flags().StringVar(&pathExact, "path-exact", "", "exact path matcher")
	cmd.Flags().StringVar(&pathPrefix, "path-prefix", "", "path prefix matcher")
	cmd.Flags().StringVar(&pathRegex, "path-regex", "", "path regex matcher")
	cmd.MarkFlagFilename("filename")
	return cmd
}

func route(event, pathExact, pathRegex, pathPrefix, filename string) (*v1.Route, error) {
	if event != "" {
		return &v1.Route{
			Matcher: &v1.Route_EventMatcher{
				EventMatcher: &v1.EventMatcher{EventType: event},
			},
		}, nil
	}

	if pathExact != "" {
		return &v1.Route{
			Matcher: &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path: &v1.RequestMatcher_PathExact{PathExact: pathExact},
				},
			},
		}, nil
	}

	if pathRegex != "" {
		return &v1.Route{
			Matcher: &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path: &v1.RequestMatcher_PathRegex{PathRegex: pathRegex},
				},
			},
		}, nil
	}

	if pathPrefix != "" {
		return &v1.Route{
			Matcher: &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path: &v1.RequestMatcher_PathPrefix{PathPrefix: pathPrefix},
				},
			},
		}, nil
	}

	return parseFile(filename)
}

func runDelete(sc storage.Interface, vhost string, route *v1.Route) ([]*v1.Route, error) {
	v, err := sc.V1().VirtualHosts().Get(vhost)
	if err != nil {
		return nil, err
	}
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
		if !match(r, route) {
			updated = append(updated, r)
		}
	}
	return updated
}

func match(left, right *v1.Route) bool {
	switch lm := left.Matcher.(type) {
	case *v1.Route_EventMatcher:
		switch rm := right.Matcher.(type) {
		case *v1.Route_EventMatcher:
			return lm.EventMatcher.EventType == rm.EventMatcher.EventType
		default:
			return false
		}
	case *v1.Route_RequestMatcher:
		switch rm := right.Matcher.(type) {
		case *v1.Route_RequestMatcher:
			return matchPath(lm.RequestMatcher, rm.RequestMatcher)
		default:
			return false
		}
	default:
		return false
	}
}

func matchPath(left, right *v1.RequestMatcher) bool {
	switch lp := left.Path.(type) {
	case *v1.RequestMatcher_PathExact:
		switch rp := right.Path.(type) {
		case *v1.RequestMatcher_PathExact:
			return lp.PathExact == rp.PathExact
		default:
			return false
		}
	case *v1.RequestMatcher_PathRegex:
		switch rp := right.Path.(type) {
		case *v1.RequestMatcher_PathRegex:
			return lp.PathRegex == rp.PathRegex
		default:
			return false
		}
	case *v1.RequestMatcher_PathPrefix:
		switch rp := right.Path.(type) {
		case *v1.RequestMatcher_PathPrefix:
			return lp.PathPrefix == rp.PathPrefix
		default:
			return false
		}
	default:
		return false
	}
}
