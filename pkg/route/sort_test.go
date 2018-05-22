package route

import (
	"bytes"
	"testing"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/protoutil"
	"github.com/solo-io/glooctl/pkg/util"
)

type TestCase [][]*v1.Route

var (
	routeExtensions, _ = protoutil.MarshalStruct(map[string]interface{}{
		"auth": map[string]interface{}{
			"credentials": struct {
				Username, Password string
			}{
				Username: "alice",
				Password: "bob",
			},
			"token": "my-12345",
		}})
	route1 = &v1.Route{
		Matcher: &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path: &v1.RequestMatcher_PathExact{
					PathExact: "/bar",
				},
				Verbs: []string{"GET", "POST"},
			},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Upstream{
				Upstream: &v1.UpstreamDestination{
					Name: "my-upstream",
				},
			},
		},
		Extensions: routeExtensions,
	}

	route2 = &v1.Route{
		Matcher: &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path: &v1.RequestMatcher_PathPrefix{
					PathPrefix: "/foo",
				},
				Headers: map[string]string{"x-foo-bar": ""},
				Verbs:   []string{"GET", "POST"},
			},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					FunctionName: "foo",
					UpstreamName: "aws",
				},
			},
		},
		Extensions: routeExtensions,
	}

	route3 = &v1.Route{
		Matcher: &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path: &v1.RequestMatcher_PathPrefix{
					PathPrefix: "/foo/bar",
				},
				Headers: map[string]string{"x-foo-bar": ""},
				Verbs:   []string{"GET", "POST"},
			},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					FunctionName: "foo",
					UpstreamName: "aws",
				},
			},
		},
		Extensions: routeExtensions,
	}

	route4 = &v1.Route{
		Matcher: &v1.Route_EventMatcher{
			EventMatcher: &v1.EventMatcher{
				EventType: "/apple",
			},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					FunctionName: "foo",
					UpstreamName: "aws",
				},
			},
		},
		Extensions: routeExtensions,
	}

	route5 = &v1.Route{
		Matcher: &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path: &v1.RequestMatcher_PathPrefix{
					PathPrefix: "/bar/foo",
				},
				Headers: map[string]string{"x-foo-bar": ""},
				Verbs:   []string{"GET", "POST"},
			},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					FunctionName: "foo",
					UpstreamName: "aws",
				},
			},
		},
		Extensions: routeExtensions,
	}
)

func TestSorting(t *testing.T) {
	data := []TestCase{
		// list 1
		{
			// unsorted
			[]*v1.Route{route1, route2, route3, route4},
			// sorted
			[]*v1.Route{route4, route1, route3, route2}},
		// list 2 - shouldn't change order if they are of similar type and length
		{
			// unsorted
			[]*v1.Route{route3, route5},
			// sorted
			[]*v1.Route{route3, route5}},
	}

	for _, tc := range data {
		SortRoutes(tc[0])
		for i, r := range tc[0] {
			if !r.Equal(tc[1][i]) {
				t.Errorf("expected %s, got %s", toString(tc[1][i]), toString(r))
			}
		}
	}
}

func toString(r *v1.Route) string {
	buf := &bytes.Buffer{}
	util.PrintYAML(r, buf)
	return buf.String()
}
