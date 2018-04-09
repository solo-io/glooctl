package vhost

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/protoutil"
)

var (
	extensions, _ = protoutil.MarshalStruct(map[string]interface{}{
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
		Extensions: extensions,
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
		Extensions: extensions,
	}
	route3 = &v1.Route{
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
		Extensions: extensions,
	}

	vhost1 = &v1.VirtualHost{
		Name:   "default",
		Routes: []*v1.Route{route1, route2, route3},
	}
)

func TestPrintYAML(t *testing.T) {
	out := captureStdOut(func() {
		printYAML(vhost1)
	})
	expected, err := ioutil.ReadFile("testdata/vhost1.yaml") // test file has extra empty line to match the print function
	if err != nil {
		t.Error("unable to load yaml file ", err)
	}

	if isDiff, e, a := diff(string(expected), out); isDiff {
		t.Errorf("expected and actual YAML didn't match:\nexpected: %s\nactual:   %s\n", e, a)
	}
}

func captureStdOut(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)
	// don't block printing
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	return <-outC
}

func diff(e, a string) (bool, string, string) {
	// doing line by line so we can show differing lines
	elines := strings.Split(e, "\n")
	alines := strings.Split(a, "\n")

	if len(elines) != len(alines) {
		return true, fmt.Sprintf("%d lines", len(elines)), fmt.Sprintf("%d lines", len(alines))
	}

	for i := range elines {
		if elines[i] != alines[i] {
			return true, elines[i], alines[i]
		}
	}
	return false, "", ""
}
