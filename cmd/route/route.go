package route

import (
	"bytes"
	"fmt"
	"strings"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo/pkg/protoutil"
	"github.com/spf13/pflag"

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
	if len(routes) == 0 {
		fmt.Println("No routes defined")
		return
	}
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

const (
	event      = "event       : "
	pathExact  = "exact path  : "
	pathRegex  = "regex path  : "
	pathPrefix = "path prefix : "
	unknown    = "matcher     : unknown"
)

func toString(r *v1.Route) string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("%s\n -> %s\n",
		matcherToString(r),
		destinationToString(r))
}

func matcherToString(r *v1.Route) string {
	switch m := r.GetMatcher().(type) {
	case *v1.Route_EventMatcher:
		return event + m.EventMatcher.EventType
	case *v1.Route_RequestMatcher:
		var path string
		switch p := m.RequestMatcher.GetPath().(type) {
		case *v1.RequestMatcher_PathExact:
			path = pathExact + p.PathExact
		case *v1.RequestMatcher_PathRegex:
			path = pathRegex + p.PathRegex
		case *v1.RequestMatcher_PathPrefix:
			path = pathPrefix + p.PathPrefix
		default:
			path = unknown
		}
		verb := ""
		if m.RequestMatcher.Verbs != nil {
			verb = fmt.Sprintf("\nmethods     : %v", m.RequestMatcher.Verbs)
		}
		headers := ""
		if m.RequestMatcher.Headers != nil {
			headers = fmt.Sprintf("\nheaders     : %v", m.RequestMatcher.Headers)
		}
		return path + verb + headers
	default:
		return unknown
	}
}

func destinationToString(r *v1.Route) string {
	single := r.GetSingleDestination()
	if single != nil {
		return upstreamToString(single.GetUpstream(), single.GetFunction())
	}

	multi := r.GetMultipleDestinations()
	if multi != nil {
		b := bytes.Buffer{}
		b.WriteString("[\n")
		for _, m := range multi {
			fmt.Fprintf(&b, "  %3d, %s\n", m.GetWeight(),
				upstreamToString(m.GetUpstream(), m.GetFunction()))
		}
		b.WriteString("]")
		return b.String()
	}

	return "unknown"
}

func upstreamToString(u *v1.UpstreamDestination, f *v1.FunctionDestination) string {
	if u != nil {
		return u.Name
	}

	if f != nil {
		return fmt.Sprintf("%s/%s", f.UpstreamName, f.FunctionName)
	}

	return "<no destintation specified>"
}

type routeDetail struct {
	event      string
	pathExact  string
	pathRegex  string
	pathPrefix string
	verb       string
	headers    string
	upstream   string
	function   string
}

func route(flags *pflag.FlagSet) (*v1.Route, error) {
	filename, _ := flags.GetString("filename")
	if filename != "" {
		return parseFile(filename)
	}

	rd := routeDetails(flags)
	return fromRouteDetail(rd)
}

func routeDetails(flags *pflag.FlagSet) *routeDetail {
	get := func(key string) string {
		v, _ := flags.GetString(key)
		return v
	}
	return &routeDetail{
		event:      get("event"),
		pathExact:  get("path-exact"),
		pathRegex:  get("path-regex"),
		pathPrefix: get("path-prefix"),
		verb:       get("http-method"),
		headers:    get("header"),
		upstream:   get("upstream"),
		function:   get("function"),
	}
}

func fromRouteDetail(rd *routeDetail) (*v1.Route, error) {
	route := &v1.Route{}

	// matcher
	if rd.event != "" {
		route.Matcher = &v1.Route_EventMatcher{
			EventMatcher: &v1.EventMatcher{EventType: rd.event},
		}
	} else {
		var verbs []string
		if rd.verb != "" {
			verbs = strings.Split(strings.ToUpper(rd.verb), ",")
			for i, v := range verbs {
				verbs[i] = strings.TrimSpace(v)
			}
		}

		var headers map[string]string
		if rd.headers != "" {
			headers = make(map[string]string)
			entries := strings.Split(rd.headers, ",")
			for _, e := range entries {
				parts := strings.SplitN(e, ":", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("unable to parse headers %s", rd.headers)
				}
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
		if rd.pathExact != "" {
			route.Matcher = &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path:    &v1.RequestMatcher_PathExact{PathExact: rd.pathExact},
					Verbs:   verbs,
					Headers: headers,
				},
			}
		} else if rd.pathRegex != "" {
			route.Matcher = &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path:    &v1.RequestMatcher_PathRegex{PathRegex: rd.pathRegex},
					Verbs:   verbs,
					Headers: headers,
				},
			}
		} else if rd.pathPrefix != "" {
			route.Matcher = &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path:    &v1.RequestMatcher_PathPrefix{PathPrefix: rd.pathPrefix},
					Verbs:   verbs,
					Headers: headers,
				},
			}
		} else {
			return nil, fmt.Errorf("a matcher wasn't specified")
		}
	}

	// destination
	if rd.upstream == "" {
		return nil, fmt.Errorf("an upstream is necessary for specifying destination")
	}
	// currently only support single destination
	if rd.function != "" {
		route.SingleDestination = &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					UpstreamName: rd.upstream,
					FunctionName: rd.function},
			},
		}
	} else if rd.upstream != "" {
		route.SingleDestination = &v1.Destination{
			DestinationType: &v1.Destination_Upstream{
				Upstream: &v1.UpstreamDestination{Name: rd.upstream},
			},
		}
	} else {
		return nil, fmt.Errorf("a destintation wasn't specified")
	}

	return route, nil
}

const defaultVHost = "default"

func createDefaultVHost(sc storage.Interface) error {
	vhost := &v1.VirtualHost{
		Name: defaultVHost,
	}
	_, err := sc.V1().VirtualHosts().Create(vhost)
	if err != nil && !storage.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func virtualHost(sc storage.Interface, name string) (*v1.VirtualHost, error) {
	// make sure default virtual host exists
	if err := createDefaultVHost(sc); err != nil {
		return nil, err
	}

	return sc.V1().VirtualHosts().Get(name)
}
