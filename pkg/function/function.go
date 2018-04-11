package function

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/olekukonko/tablewriter"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/protoutil"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/route"
)

type FunctionWithUpstream struct {
	Upstream string
	Function *v1.Function
}

func Get(sc storage.Interface, output, tplt string) error {
	upstreams, err := sc.V1().Upstreams().List()
	if err != nil {
		return err
	}
	if len(upstreams) == 0 {
		fmt.Println("No upstreams found.")
		return nil
	}

	var functions []FunctionWithUpstream
	for _, u := range upstreams {
		funcs := u.GetFunctions()
		for _, f := range funcs {
			functions = append(functions, FunctionWithUpstream{
				Upstream: u.Name,
				Function: f,
			})
		}
	}
	switch output {
	case "yaml":
		return printYAMLList(functions)
	case "json":
		return printJSONList(functions)
	case "template":
		return PrintTemplate(functions, tplt, os.Stdout)
	default:
		virtualhosts, err := sc.V1().VirtualHosts().List()
		if err != nil {
			return errors.Wrap(err, "unable to get virtual hosts")
		}
		PrintTableWithRoutes(functions, os.Stdout, virtualhosts)
	}
	return nil
}

func printYAMLList(list []FunctionWithUpstream) error {
	for _, f := range list {
		jsn, err := protoutil.Marshal(f.Function)
		if err != nil {
			return errors.Wrap(err, "unable to marshal")
		}
		b, err := yaml.JSONToYAML(jsn)
		if err != nil {
			return errors.Wrap(err, "unable to convert to YAML")
		}
		fmt.Println(string(b))
	}
	return nil
}

func printJSONList(list []FunctionWithUpstream) error {
	for _, f := range list {
		b, err := protoutil.Marshal(f.Function)
		if err != nil {
			return errors.Wrap(err, "unable to conver to JSON")
		}
		fmt.Println(string(b))
	}
	return nil
}

func PrintTemplate(list []FunctionWithUpstream, tplt string, w io.Writer) error {
	t, err := template.New("output").Parse(tplt)
	if err != nil {
		return errors.Wrap(err, "unable to parse template")
	}
	return t.Execute(w, list)
}

func PrintTableWithRoutes(list []FunctionWithUpstream, w io.Writer, virtualhosts []*v1.VirtualHost) {
	routeMap := routes(virtualhosts)
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Function", "Upstream", "Matcher", "Type", "Verb", "Header", "Extension"})

	for _, f := range list {
		key := fmt.Sprintf("%s:%s", f.Upstream, f.Function.Name)
		routes, ok := routeMap[key]
		if !ok {
			table.Append([]string{f.Function.Name, f.Upstream, "", "", "", "", ""})
		} else {
			for _, r := range routes {
				matcher, rType, verb, headers := route.Matcher(r)
				ext := route.Extension(r)
				table.Append([]string{f.Function.Name, f.Upstream,
					matcher, rType, verb, headers, ext})
			}
		}
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func routes(virtualhosts []*v1.VirtualHost) map[string][]*v1.Route {
	m := make(map[string][]*v1.Route)
	for _, v := range virtualhosts {
		for _, r := range v.Routes {
			dsts := route.Destinations(r)
			for _, d := range dsts {
				if d.Function != "" {
					key := fmt.Sprintf("%s:%s", d.Upstream, d.Function)
					m[key] = append(m[key], r)
				}
			}
		}
	}
	return m
}
