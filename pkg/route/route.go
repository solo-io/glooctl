package route

import (
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/solo-io/gloo/pkg/api/types/v1"
)

type Destination struct {
	Upstream string
	Function string
}

func PrintTable(list []*v1.Route, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Matcher", "Type", "Verb", "Header", "Upstream", "Function", "Extension"})

	for _, r := range list {
		matcher, rType, verb, headers := Matcher(r)
		ext := Extension(r)
		for _, d := range Destinations(r) {
			table.Append([]string{matcher, rType, verb, headers, d.Upstream, d.Function, ext})
		}
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func Matcher(r *v1.Route) (string, string, string, string) {
	switch m := r.GetMatcher().(type) {
	case *v1.Route_EventMatcher:
		return m.EventMatcher.EventType, "Event", "", ""
	case *v1.Route_RequestMatcher:
		var path string
		var rType string
		switch p := m.RequestMatcher.GetPath().(type) {
		case *v1.RequestMatcher_PathExact:
			path = p.PathExact
			rType = "Exact Path"
		case *v1.RequestMatcher_PathPrefix:
			path = p.PathPrefix
			rType = "Path Prefix"
		case *v1.RequestMatcher_PathRegex:
			path = p.PathRegex
			rType = "Regex Path"
		default:
			path = ""
			rType = "Unknown"
		}
		verb := "*"
		if m.RequestMatcher.Verbs != nil {
			verb = strings.Join(m.RequestMatcher.Verbs, " ")
		}
		headers := ""
		if m.RequestMatcher.Headers != nil {
			builder := strings.Builder{}
			for k, v := range m.RequestMatcher.Headers {
				builder.WriteString(k)
				builder.WriteString(":")
				builder.WriteString(v)
				builder.WriteString("; ")
			}
			headers = builder.String()
		}
		return path, rType, verb, headers
	default:
		return "", "Unknown", "", ""
	}
}

func Destinations(r *v1.Route) []Destination {
	single := r.GetSingleDestination()
	if single != nil {
		return []Destination{upstreamToDestination(single.GetUpstream(), single.GetFunction())}
	}

	multi := r.GetMultipleDestinations()
	if multi != nil {
		d := make([]Destination, len(multi))
		for i, m := range multi {
			d[i] = upstreamToDestination(m.GetUpstream(), m.GetFunction())
		}
		return d
	}

	return []Destination{Destination{"", ""}}
}

func upstreamToDestination(u *v1.UpstreamDestination, f *v1.FunctionDestination) Destination {
	if u != nil {
		return Destination{u.Name, ""}
	}

	if f != nil {
		return Destination{f.UpstreamName, f.FunctionName}
	}

	return Destination{"", ""}
}

func Extension(r *v1.Route) string {
	ext := r.GetExtensions()
	if ext == nil || ext.GetFields() == nil {
		return ""
	}

	builder := strings.Builder{}
	for k, v := range ext.GetFields() {
		builder.WriteString(k)
		builder.WriteString(":")
		builder.WriteString(v.GoString())
		builder.WriteString("; ")
	}
	return builder.String()
}
