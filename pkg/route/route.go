package route

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/gogo/protobuf/types"
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
	s := make([]string, len(ext.GetFields()))
	i := 0
	for k, v := range ext.GetFields() {
		s[i] = fmt.Sprintf("%s: %s", k, prettyPrint(v))
		i++
	}
	return fmt.Sprintf("{%s}", strings.Join(s, ", "))
}

func prettyPrint(v *types.Value) string {
	switch t := v.Kind.(type) {
	case *types.Value_NullValue:
		return ""
	case *types.Value_NumberValue:
		return fmt.Sprintf("%v", t.NumberValue)
	case *types.Value_StringValue:
		return fmt.Sprintf("\"%v\"", t.StringValue)
	case *types.Value_BoolValue:
		return fmt.Sprintf("%v", t.BoolValue)
	case *types.Value_StructValue:
		return prettyPrintStruct(t)
	case *types.Value_ListValue:
		return prettyPrintList(t)
	default:
		return "<unknown>"
	}
}

func prettyPrintList(lv *types.Value_ListValue) string {
	if lv == nil || lv.ListValue == nil || lv.ListValue.Values == nil {
		return ""
	}
	s := make([]string, len(lv.ListValue.Values))
	for i, v := range lv.ListValue.Values {
		s[i] = prettyPrint(v)
	}
	return fmt.Sprintf("[%s]", strings.Join(s, ", "))
}

func prettyPrintStruct(sv *types.Value_StructValue) string {
	if sv == nil || sv.StructValue == nil || sv.StructValue.Fields == nil {
		return ""
	}

	s := make([]string, len(sv.StructValue.GetFields()))
	i := 0
	for k, v := range sv.StructValue.GetFields() {
		s[i] = fmt.Sprintf("%s: %s", k, prettyPrint(v))
		i++
	}
	return fmt.Sprintf("{%s}", strings.Join(s, ", "))

}
