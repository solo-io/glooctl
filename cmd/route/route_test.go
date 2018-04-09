package route

import (
	"fmt"
	"testing"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/spf13/cobra"
)

type routeDetailTestCase struct {
	args  []string
	route *v1.Route
}

func TestToRoute(t *testing.T) {
	cases := []routeDetailTestCase{
		routeDetailTestCase{
			args: []string{"--event", "apple", "--upstream", "aws", "--function", "func1"},
			route: &v1.Route{
				Matcher: &v1.Route_EventMatcher{
					EventMatcher: &v1.EventMatcher{
						EventType: "apple",
					},
				},
				SingleDestination: &v1.Destination{
					DestinationType: &v1.Destination_Function{
						Function: &v1.FunctionDestination{
							UpstreamName: "aws",
							FunctionName: "func1",
						},
					},
				},
			},
		},
		routeDetailTestCase{
			args: []string{"--path-prefix", "/foo", "--upstream", "foo", "--http-method", "get,put, post"},
			route: &v1.Route{
				Matcher: &v1.Route_RequestMatcher{
					RequestMatcher: &v1.RequestMatcher{
						Path:  &v1.RequestMatcher_PathPrefix{PathPrefix: "/foo"},
						Verbs: []string{"GET", "PUT", "POST"},
					},
				},
				SingleDestination: &v1.Destination{
					DestinationType: &v1.Destination_Upstream{
						Upstream: &v1.UpstreamDestination{
							Name: "foo",
						},
					},
				},
			},
		},
		routeDetailTestCase{
			args: []string{"--path-exact", "/foo", "--upstream", "foo", "--header", "key1:value1, key2:value2"},
			route: &v1.Route{
				Matcher: &v1.Route_RequestMatcher{
					RequestMatcher: &v1.RequestMatcher{
						Path:    &v1.RequestMatcher_PathExact{PathExact: "/foo"},
						Headers: map[string]string{"key1": "value1", "key2": "value2"},
					},
				},
				SingleDestination: &v1.Destination{
					DestinationType: &v1.Destination_Upstream{
						Upstream: &v1.UpstreamDestination{
							Name: "foo",
						},
					},
				},
			},
		},
	}

	for i, tc := range cases {
		cmd := &cobra.Command{}
		setupRouteParams(cmd)
		flags := cmd.Flags()
		flags.Parse(tc.args)
		route, err := fromRouteDetail(routeDetails(flags))
		if err != nil {
			t.Errorf("case %d failed conversion", i)
		}
		if !tc.route.Equal(route) {
			fmt.Println("expected: ")
			printYAML(tc.route)
			fmt.Println("got: ")
			printYAML(route)
			t.Errorf("case %d failed comparison", i)
		}
	}
}
