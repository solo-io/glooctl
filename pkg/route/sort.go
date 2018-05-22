package route

import (
	"sort"

	"github.com/solo-io/gloo/pkg/api/types/v1"
)

func SortRoutes(routes []*v1.Route) {
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
			return lessRequestMatcher(l.RequestMatcher, r.RequestMatcher)
		}
	}

	return true
}

func lessRequestMatcher(left, right *v1.RequestMatcher) bool {
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
