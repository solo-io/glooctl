package virtualservice

import (
	"io"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/glooctl/pkg/route"
)

// PrintTemplate prints virtual services using the provided Go template to io.Writer
func PrintTemplate(list []*v1.VirtualService, tmpl string, w io.Writer) error {
	t, err := template.New("output").Parse(tmpl)
	if err != nil {
		return errors.Wrap(err, "unable to parse template")
	}
	return t.Execute(w, list)
}

// PrintTable prints virtual services using tables to io.Writer
func PrintTable(list []*v1.VirtualService, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Virtual Service", "Domain", "SSL", "Status", "Matcher", "Type", "Verb", "Header", "Upstream", "Function", "Extension"})

	for _, v := range list {
		name := v.GetName()
		d := domains(v)
		ssl := sslConfig(v)
		s := status(v)

		if v.GetRoutes() == nil || len(v.GetRoutes()) == 0 {
			table.Append([]string{name, d, ssl, s, "", "", "", "", "", "", ""})
		} else {
			for i, r := range v.GetRoutes() {
				matcher, rType, verb, headers := route.Matcher(r)
				ext := route.Extension(r)
				for _, dst := range route.Destinations(r) {
					if i == 0 {
						table.Append([]string{name, d, ssl, s, matcher, rType, verb, headers, dst.Upstream, dst.Function, ext})
					} else {

						table.Append([]string{"", "", "", "", matcher, rType, verb, headers, dst.Upstream, dst.Function, ext})
					}
				}
			}
		}
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func domains(v *v1.VirtualService) string {
	if v.GetDomains() == nil || len(v.GetDomains()) == 0 {
		return ""
	}

	return strings.Join(v.GetDomains(), ", ")
}

func sslConfig(v *v1.VirtualService) string {
	if v.GetSslConfig() == nil {
		return ""
	}
	return v.GetSslConfig().GetSecretRef()
}

func status(v *v1.VirtualService) string {
	if v.Status == nil {
		return ""
	}
	return v.Status.State.String()
}
