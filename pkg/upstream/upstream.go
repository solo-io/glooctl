package upstream

import (
	"github.com/olekukonko/tablewriter"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"io"
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
			for _, f := range u.Functions {
				table.Append([]string{u.Name, u.Type, status, f.Name})
			}
		} else {
			table.Append([]string{name, uType, status, ""})
		}
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}
