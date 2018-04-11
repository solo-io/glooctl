package secret

import (
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	secret "github.com/solo-io/gloo-secret"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugins/aws"
	"github.com/solo-io/gloo/pkg/plugins/google"
)

const (
	googleAnnotationKey = "gloo.solo.io/google_secret_ref"
)

func PrintTableWithUsage(list []*secret.Secret, w io.Writer, u []*v1.Upstream, v []*v1.VirtualHost) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Type", "In Use By"})

	usageMap := usage(list, u, v)
	for _, s := range list {
		name := s.Name
		sType := secretType(s)
		usage := ""
		use, ok := usageMap[name]
		if ok {
			usage = strings.Join(use, "; ")
		}
		table.Append([]string{name, sType, usage})
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func usage(list []*secret.Secret, upstreams []*v1.Upstream, virtualhosts []*v1.VirtualHost) map[string][]string {
	m := make(map[string][]string, len(list))
	for _, s := range list {
		m[s.Name] = []string{}
	}

	if virtualhosts != nil {
		for _, v := range virtualhosts {
			ssl := v.GetSslConfig()
			if ssl != nil {
				ref := ssl.GetSecretRef()
				existing, ok := m[ref]
				if ok {
					m[ref] = append(existing, "Virtual Host:"+v.GetName())
				}
			}
		}
	}

	if upstreams != nil {
		for _, u := range upstreams {
			switch u.Type {
			case aws.UpstreamTypeAws:
				spec, err := aws.DecodeUpstreamSpec(u.Spec)
				if err != nil {
					continue // TODO log it
				}
				ref := spec.SecretRef
				existing, ok := m[ref]
				if ok {
					m[ref] = append(existing, "Upstream:"+u.Name)
				}
			case gfunc.UpstreamTypeGoogle:
				if u.Metadata == nil || u.Metadata.Annotations == nil {
					continue
				}
				ref, ok := u.Metadata.Annotations[googleAnnotationKey]
				if !ok {
					continue
				}
				existing, ok := m[ref]
				if ok {
					m[ref] = append(existing, "Upstream:"+u.Name)
				}
			}
		}
	}

	return m
}

func secretType(s *secret.Secret) string {
	if s.Data == nil {
		return "Unknown"
	}
	if _, ok := s.Data[serviceAccountJsonKeyFile]; ok {
		return "Google"
	}

	_, first := s.Data[awsAccessKey]
	_, second := s.Data[awsSecretKey]
	if first && second {
		return "AWS"
	}

	_, first = s.Data[sslCertificateChainKey]
	_, second = s.Data[sslPrivateKeyKey]
	if first && second {
		return "Certificate"
	}
	return "Unknown"
}
