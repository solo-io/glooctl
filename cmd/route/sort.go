package route

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
	proute "github.com/solo-io/glooctl/pkg/route"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func sortCmd(opts *client.StorageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sort",
		Short: "sort routes to have the longest route first",
		Run: func(c *cobra.Command, args []string) {
			sc, err := client.StorageClient(opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			domain, _ := c.InheritedFlags().GetString("domain")
			vhostname, _ := c.InheritedFlags().GetString(flagVirtualHost)

			routes, err := runSort(sc, vhostname, domain)
			if err != nil {

			}
			output, _ := c.InheritedFlags().GetString("output")
			util.PrintList(output, "", routes,
				func(data interface{}, w io.Writer) error {
					proute.PrintTable(data.([]*v1.Route), w)
					return nil
				}, os.Stdout)
		},
	}
	return cmd
}

func runSort(sc storage.Interface, vhostname, domain string) ([]*v1.Route, error) {
	v, err := virtualHost(sc, vhostname, domain, false)
	if err != nil {
		return nil, err
	}
	fmt.Println("Using virtual host:", v.Name)
	sortRoutes(v.Routes)
	updated, err := sc.V1().VirtualHosts().Update(v)
	if err != nil {
		return nil, err
	}
	return updated.GetRoutes(), nil
}

func sortRoutes(routes []*v1.Route) {
	sort.SliceStable(routes, func(i, j int) bool {
		return lessRoutes(routes[i], routes[j])
	})
}

func lessRoutes(left, right *v1.Route) bool {
	lm := left.GetMatcher()
	rm := right.GetMatcher()

	switch l := lm.(type) {
	case *v1.Route_EventMatcher:
		switch r := rm.(type) {
		case *v1.Route_EventMatcher:
			return len(l.EventMatcher.EventType) > len(r.EventMatcher.EventType)
		case *v1.Route_RequestMatcher:
			return true
		}
	case *v1.Route_RequestMatcher:
		switch r := rm.(type) {
		case *v1.Route_EventMatcher:
			return false
		case *v1.Route_RequestMatcher:
			return less(l.RequestMatcher, r.RequestMatcher)
		}
	}

	return true
}

func less(left, right *v1.RequestMatcher) bool {
	lp := left.GetPath()
	rp := right.GetPath()

	switch l := lp.(type) {
	case *v1.RequestMatcher_PathExact:
		switch r := rp.(type) {
		case *v1.RequestMatcher_PathExact:
			return len(l.PathExact) > len(r.PathExact)
		case *v1.RequestMatcher_PathRegex:
			return true
		case *v1.RequestMatcher_PathPrefix:
			return true
		}
	case *v1.RequestMatcher_PathRegex:
		switch r := rp.(type) {
		case *v1.RequestMatcher_PathExact:
			return false
		case *v1.RequestMatcher_PathRegex:
			return len(l.PathRegex) > len(r.PathRegex)
		case *v1.RequestMatcher_PathPrefix:
			return true
		}
	case *v1.RequestMatcher_PathPrefix:
		switch r := rp.(type) {
		case *v1.RequestMatcher_PathExact:
			return false
		case *v1.RequestMatcher_PathRegex:
			return false
		case *v1.RequestMatcher_PathPrefix:
			return len(l.PathPrefix) > len(r.PathPrefix)
		}
	}

	return true
}
