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
	flagDomain      = "domain"
	flagVirtualHost = "virtual-host"
	flagFilename    = "filename"

	flagEvent         = "event"
	flagPathExact     = "path-exact"
	flagPathRegex     = "path-regex"
	flagPathPrefix    = "path-prefix"
	flagMethod        = "http-method"
	flagHeaders       = "header"
	flagUpstream      = "upstream"
	flagFunction      = "function"
	flagPrefixRewrite = "prefix-rewrite"
	flagExtension     = "extensions"

	flagKubeName      = "kube-upstream"
	flagKubeNamespace = "kube-namespace"
	flagKubePort      = "kube-port"

	defaultVHost = "default"

	upstreamTypeKubernetes = "kubernetes"
	kubeSpecName           = "service_name"
	kubeSpecNamespace      = "service_namespace"
	kubeSpecPort           = "service_port"
)

type routeOption struct {
	route       *routeDetail
	virtualhost string
	domain      string
	filename    string
	output      string
	sort        bool
	interactive bool
}

type kubeUpstream struct {
	name      string
	namespace string
	port      int
}

type routeDetail struct {
	event         string
	pathExact     string
	pathRegex     string
	pathPrefix    string
	verb          string
	headers       string
	upstream      string
	function      string
	prefixRewrite string
	extensions    string

	kube *kubeUpstream
}

func parseFile(filename string) (*v1.Route, error) {
	var r v1.Route
	err := file.ReadFileInto(filename, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func route(opts *routeOption, sc storage.Interface) (*v1.Route, error) {
	if opts.filename != "" {
		return parseFile(opts.filename)
	}

	rd := opts.route
	if rd.kube.name != "" {
		upstream, err := upstream(rd.kube, sc)
		if err != nil {
			return nil, err
		}
		rd.upstream = upstream.Name
	}
	return fromRouteDetail(rd)
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

	// prefix rewrite
	if rd.prefixRewrite != "" {
		route.PrefixRewrite = rd.prefixRewrite
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

	// extensions
	if rd.extensions != "" {
		ext := &google_protobuf.Struct{}

		// special case: reading from stdin
		if rd.extensions == "-" {
			if err := util.ReadStdinInto(ext); err != nil {
				return nil, errors.Wrap(err, "reading extensions from stdin")
			}
		} else {
			err := file.ReadFileInto(rd.extensions, ext)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to read file %s for extensions", rd.extensions)
			}
		}
		route.Extensions = ext

	}
	return route, nil
}

func upstream(kube *kubeUpstream, sc storage.Interface) (*v1.Upstream, error) {
	upstreams, err := sc.V1().Upstreams().List()
	if err != nil {
		return nil, err
	}
	for _, u := range upstreams {
		if u.Type != upstreamTypeKubernetes {
			continue
		}
		s, _ := protoutil.MarshalMap(u.Spec)
		n, exists := s[kubeSpecName].(string)
		if !exists {
			continue
		}
		if n != kube.name {
			continue
		}
		if kube.namespace != "" {
			ns, exists := s[kubeSpecNamespace].(string)
			if !exists {
				continue
			}
			if ns != kube.namespace {
				continue
			}
		}

		if kube.port != 0 {
			p, exists := s[kubeSpecPort].(string)
			if !exists {
				continue
			}
			if p != strconv.Itoa(kube.port) {
				continue
			}
		}
		return u, nil
	}
	return nil, fmt.Errorf("unable to find kubernetes upstream %s/%s", kube.namespace, kube.name)
}
