package upstream

import (
	"io"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
)

func PrintTable(list []*v1.Upstream, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Type", "Status", "Function"})

	for _, u := range list {
		name := u.Name
		uType := u.Type
		if u.ServiceInfo != nil && u.ServiceInfo.Type != "" {
			uType = u.ServiceInfo.Type
		}
		status := ""
		if u.Status != nil {
			status = u.Status.State.String()
		}

		if u.Functions != nil && len(u.Functions) > 0 {
			for i, f := range u.Functions {
				if i == 0 {
					table.Append([]string{u.Name, u.Type, status, f.Name})
				} else {
					table.Append([]string{"", "", "", f.Name})
				}
			}
		} else {
			table.Append([]string{name, uType, status, ""})
		}
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func PrintTemplate(list []*v1.Upstream, tmpl string, w io.Writer) error {
	t, err := template.New("output").Parse(tmpl)
	if err != nil {
		return errors.Wrap(err, "unable to parse template")
	}
	return t.Execute(w, list)
}
