package route

import (
	"bytes"
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

type Option struct {
	Route          *Detail
	Virtualservice string
	Domain         string
	Filename       string
	Output         string
	Sort           bool
	Interactive    bool
	Index          int
}

type KubeUpstream struct {
	Name      string
	Namespace string
	Port      int
}

type Detail struct {
	Event            string
	PathExact        string
	PathRegex        string
	PathPrefix       string
	Verb             string
	Headers          string
	Upstream         string
	Function         string
	PrefixRewrite    string
	Extensions       string
	InlineExtensions string

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

func FromOption(opts *Option, sc storage.Interface) (*v1.Route, error) {
	if opts.Filename != "" {
		return parseFile(opts.Filename)
	}

	if opts.Index > 0 {
		return fromIndex(opts.Index, opts.Virtualservice, sc)
	}

	rd := opts.Route
	if rd.Kube.Name != "" {
		upstream, err := upstream(rd.Kube, sc)
		if err != nil {
			return nil, err
		}
		rd.Upstream = upstream.Name
	}

	return FromDetail(rd)
}

func fromIndex(index int, virtualService string, sc storage.Interface) (*v1.Route, error) {
	vs, err := sc.V1().VirtualServices().Get(virtualService)
	if err != nil {
		return nil, err
	}
	if len(vs.Routes) < index {
		return nil, errors.Errorf("invalid index %d; should be between 1 and %d", index, len(vs.Routes))
	}

	return vs.Routes[index-1], nil
}

func FromDetail(rd *Detail) (*v1.Route, error) {
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

func destinationFromDetails(rd *Detail) (*v1.Destination, error) {
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

func extensionsFromDetails(rd *Detail) (*google_protobuf.Struct, error) {
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

func ToDetails(r *v1.Route) (*Detail, error) {
	rd := &Detail{}

	// matcher
	switch m := r.GetMatcher().(type) {
	case *v1.Route_EventMatcher:
		rd.Event = m.EventMatcher.EventType
	case *v1.Route_RequestMatcher:
		switch p := m.RequestMatcher.GetPath().(type) {
		case *v1.RequestMatcher_PathExact:
			rd.PathExact = p.PathExact
		case *v1.RequestMatcher_PathPrefix:
			rd.PathPrefix = p.PathPrefix
		case *v1.RequestMatcher_PathRegex:
			rd.PathRegex = p.PathRegex
		}

		if len(m.RequestMatcher.Verbs) > 0 {
			rd.Verb = strings.Join(m.RequestMatcher.Verbs, ",")
		}
		builder := bytes.Buffer{}
		for k, v := range m.RequestMatcher.Headers {
			builder.WriteString(k)
			builder.WriteString(":")
			builder.WriteString(v)
			builder.WriteString("; ")
		}
		rd.Headers = builder.String()
	}

	// destination
	dstList := Destinations(r)
	// use the first one; TODO support for multiple desitinations
	rd.Upstream = dstList[0].Upstream
	rd.Function = dstList[0].Function

	// additional
	rd.PrefixRewrite = r.PrefixRewrite
	if r.Extensions != nil {
		b, err := protoutil.Marshal(r.Extensions)
		if err != nil {
			return nil, err
		}
		rd.InlineExtensions = string(b)
	}

	return rd, nil
}
