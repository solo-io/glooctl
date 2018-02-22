package route

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/protoutil"

	"github.com/ghodss/yaml"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage/file"
)

func parseFile(filename string) (*v1.Route, error) {
	var r v1.Route
	err := file.ReadFileInto(filename, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func printRoutes(routes []*v1.Route, output string) {
	switch output {
	case "json":
		printJSONList(routes)
	case "yaml":
		printYAMLList(routes)
	default:
		printSummaryList(routes)
	}
}

func printJSON(r *v1.Route) {
	b, err := protoutil.Marshal(r)
	if err != nil {
		fmt.Println("unable to convert to JSON ", err)
		return
	}
	fmt.Println(string(b))
}

func printYAML(r *v1.Route) {
	jsn, err := protoutil.Marshal(r)
	if err != nil {
		fmt.Println("unable to marshal ", err)
		return
	}
	b, err := yaml.JSONToYAML(jsn)
	if err != nil {
		fmt.Println("unable to convert to YAML ", err)
		return
	}
	fmt.Println(string(b))
}

func printJSONList(routes []*v1.Route) {
	for _, r := range routes {
		printJSON(r)
	}
}

func printYAMLList(routes []*v1.Route) {
	for _, r := range routes {
		printYAML(r)
	}
}

func printSummaryList(r []*v1.Route) {
	for _, entry := range r {
		fmt.Println(toString(entry))
	}
}

func toString(r *v1.Route) string {
	if r == nil {
		return ""
	}
	switch m := r.GetMatcher().(type) {
	case *v1.Route_EventMatcher:
		return "event matcher: " + m.EventMatcher.EventType
	case *v1.Route_RequestMatcher:
		path := m.RequestMatcher.GetPath()
		switch p := path.(type) {
		case *v1.RequestMatcher_PathExact:
			return "request exact path: " + p.PathExact
		case *v1.RequestMatcher_PathRegex:
			return "request path regex: " + p.PathRegex
		case *v1.RequestMatcher_PathPrefix:
			return "request path prefix: " + p.PathPrefix
		default:
			return "route: <unknown request matcher>"
		}
	default:
		return "route: <unknown matcher>"
	}
}
