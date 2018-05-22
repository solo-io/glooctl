package route

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/protoutil"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/util"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage/file"

	_ "github.com/gogo/protobuf/gogoproto"
	google_protobuf "github.com/gogo/protobuf/types"
)

const (
	upstreamTypeKubernetes = "kubernetes"
	kubeSpecName           = "service_name"
	kubeSpecNamespace      = "service_namespace"
	kubeSpecPort           = "service_port"
)

type RouteOption struct {
	Route          *RouteDetail
	Virtualservice string
	Domain         string
	Filename       string
	Output         string
	Sort           bool
	Interactive    bool
}

type KubeUpstream struct {
	Name      string
	Namespace string
	Port      int
}

type RouteDetail struct {
	Event         string
	PathExact     string
	PathRegex     string
	PathPrefix    string
	Verb          string
	Headers       string
	Upstream      string
	Function      string
	PrefixRewrite string
	Extensions    string

	Kube *KubeUpstream
}

func parseFile(filename string) (*v1.Route, error) {
	var r v1.Route
	err := file.ReadFileInto(filename, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func FromRouteOption(opts *RouteOption, sc storage.Interface) (*v1.Route, error) {
	if opts.Filename != "" {
		return parseFile(opts.Filename)
	}

	rd := opts.Route
	if rd.Kube.Name != "" {
		upstream, err := upstream(rd.Kube, sc)
		if err != nil {
			return nil, err
		}
		rd.Upstream = upstream.Name
	}
	return FromRouteDetail(rd)
}

func FromRouteDetail(rd *RouteDetail) (*v1.Route, error) {
	route := &v1.Route{}

	// matcher
	if rd.Event != "" {
		route.Matcher = &v1.Route_EventMatcher{
			EventMatcher: &v1.EventMatcher{EventType: rd.Event},
		}
	} else {
		var verbs []string
		if rd.Verb != "" {
			verbs = strings.Split(strings.ToUpper(rd.Verb), ",")
			for i, v := range verbs {
				verbs[i] = strings.TrimSpace(v)
			}
		}

		var headers map[string]string
		if rd.Headers != "" {
			headers = make(map[string]string)
			entries := strings.Split(rd.Headers, ",")
			for _, e := range entries {
				parts := strings.SplitN(e, ":", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("unable to parse headers %s", rd.Headers)
				}
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
		if rd.PathExact != "" {
			route.Matcher = &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path:    &v1.RequestMatcher_PathExact{PathExact: rd.PathExact},
					Verbs:   verbs,
					Headers: headers,
				},
			}
		} else if rd.PathRegex != "" {
			route.Matcher = &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path:    &v1.RequestMatcher_PathRegex{PathRegex: rd.PathRegex},
					Verbs:   verbs,
					Headers: headers,
				},
			}
		} else if rd.PathPrefix != "" {
			route.Matcher = &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path:    &v1.RequestMatcher_PathPrefix{PathPrefix: rd.PathPrefix},
					Verbs:   verbs,
					Headers: headers,
				},
			}
		} else {
			return nil, fmt.Errorf("a matcher wasn't specified")
		}
	}

	// prefix rewrite
	if rd.PrefixRewrite != "" {
		route.PrefixRewrite = rd.PrefixRewrite
	}

	// destination
	dst, err := destinationFromDetails(rd)
	if err != nil {
		return nil, err
	}
	// currently only support single destination from CLI
	route.SingleDestination = dst

	// extensions
	ext, err := extensionsFromDetails(rd)
	if err != nil {
		return nil, err
	}
	if ext != nil {
		route.Extensions = ext
	}
	return route, nil
}

func destinationFromDetails(rd *RouteDetail) (*v1.Destination, error) {
	if rd.Upstream == "" {
		return nil, fmt.Errorf("an upstream is necessary for specifying destination")
	}
	// currently only support single destination
	if rd.Function != "" {
		return &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					UpstreamName: rd.Upstream,
					FunctionName: rd.Function},
			},
		}, nil
	}

	return &v1.Destination{
		DestinationType: &v1.Destination_Upstream{
			Upstream: &v1.UpstreamDestination{Name: rd.Upstream},
		},
	}, nil
}

func extensionsFromDetails(rd *RouteDetail) (*google_protobuf.Struct, error) {
	if rd.Extensions == "" {
		return nil, nil
	}
	ext := &google_protobuf.Struct{}

	// special case: reading from stdin
	if rd.Extensions == "-" {
		if err := util.ReadStdinInto(ext); err != nil {
			return nil, errors.Wrap(err, "reading extensions from stdin")
		}
		return ext, nil
	}

	err := file.ReadFileInto(rd.Extensions, ext)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file %s for extensions", rd.Extensions)
	}

	return ext, nil
}

func upstream(kube *KubeUpstream, sc storage.Interface) (*v1.Upstream, error) {
	upstreams, err := sc.V1().Upstreams().List()
	if err != nil {
		return nil, err
	}
	for _, u := range upstreams {
		if u.Type != upstreamTypeKubernetes {
			continue
		}
		s, err := protoutil.MarshalMap(u.Spec)
		if err != nil {
			return nil, err
		}
		n, exists := s[kubeSpecName].(string)
		if !exists {
			continue
		}
		if n != kube.Name {
			continue
		}
		if kube.Namespace != "" {
			ns, exists := s[kubeSpecNamespace].(string)
			if !exists {
				continue
			}
			if ns != kube.Namespace {
				continue
			}
		}

		if kube.Port != 0 {
			p, exists := s[kubeSpecPort].(string)
			if !exists {
				continue
			}
			if p != strconv.Itoa(kube.Port) {
				continue
			}
		}
		return u, nil
	}
	return nil, fmt.Errorf("unable to find kubernetes upstream %s/%s", kube.Namespace, kube.Name)
}
