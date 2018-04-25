package function

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/olekukonko/tablewriter"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/route"
	"github.com/solo-io/glooctl/pkg/util"
)

type FunctionWithUpstream struct {
	Upstream string
	Function *v1.Function
}

// Get gets a list of functions defined in the system and prints them
// in the format specified by output
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
		return util.PrintYAMLList(toV1Functions(functions), os.Stdout)
	case "json":
		return util.PrintJSONList(toV1Functions(functions), os.Stdout)
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

func toV1Functions(list []FunctionWithUpstream) []*v1.Function {
	functions := make([]*v1.Function, len(list))
	for i, f := range list {
		functions[i] = f.Function
	}
	return functions
}

// PrintTemplate prints functions using the provided Go template to the io.Writer
func PrintTemplate(list []FunctionWithUpstream, tplt string, w io.Writer) error {
	t, err := template.New("output").Parse(tplt)
	if err != nil {
		return errors.Wrap(err, "unable to parse template")
	}
	return t.Execute(w, list)
}

// PrintTableWithRoutes prints functions and routes mapped to them to the io.Writer
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
