package secret

import (
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugins/aws"
	"github.com/solo-io/gloo/pkg/plugins/google"
	"github.com/solo-io/gloo/pkg/storage/dependencies"
)

const (
	// GoogleAnnotationKey is the key for annotation used in Google upstream
	GoogleAnnotationKey = "gloo.solo.io/google_secret_ref"
)

type upstreamSecretRefFetcher func(*v1.Upstream) []string

var (
	upstreamUsagePlugins = map[string]upstreamSecretRefFetcher{
		aws.UpstreamTypeAws:      checkAWS,
		gfunc.UpstreamTypeGoogle: checkGoogle,
	}
)

// PrintTableWithUsage prints secrets and their usage
func PrintTableWithUsage(list []*dependencies.Secret, w io.Writer, u []*v1.Upstream, v []*v1.VirtualService) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Type", "In Use By"})

	usageMap := usage(list, u, v)
	for _, s := range list {
		name := s.Ref
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

func usage(list []*dependencies.Secret, upstreams []*v1.Upstream, virtualservices []*v1.VirtualService) map[string][]string {
	m := make(map[string][]string, len(list))
	for _, s := range list {
		m[s.Ref] = []string{}
	}

	usageVirtualService(m, virtualservices)
	usageUpstream(m, upstreams)
	return m
}

func usageVirtualService(secrets map[string][]string, virtualservices []*v1.VirtualService) {
	for _, v := range virtualservices {
		ssl := v.GetSslConfig()
		if ssl != nil {
			ref := ssl.GetSecretRef()
			existing, ok := secrets[ref]
			if ok {
				secrets[ref] = append(existing, "Virtual Service:"+v.GetName())
			}
		}
	}
}

func usageUpstream(secrets map[string][]string, upstreams []*v1.Upstream) {
	for _, u := range upstreams {
		fetcher, hasPlugin := upstreamUsagePlugins[u.Type]
		if !hasPlugin {
			continue
		}
		for _, ref := range fetcher(u) {
			existing, knownSecret := secrets[ref]
			if knownSecret {
				secrets[ref] = append(existing, "Upstream:"+u.Name)
			}
		}
	}
}

func checkAWS(u *v1.Upstream) []string {
	spec, err := aws.DecodeUpstreamSpec(u.Spec)
	if err != nil {
		return nil
	}
	return []string{spec.SecretRef}
}

func checkGoogle(u *v1.Upstream) []string {
	if u.Metadata == nil || u.Metadata.Annotations == nil {
		return nil
	}
	ref, ok := u.Metadata.Annotations[GoogleAnnotationKey]
	if !ok {
		return nil
	}
	return []string{ref}
}

func secretType(s *dependencies.Secret) string {
	if s.Data == nil {
		return "Unknown"
	}
	if _, ok := s.Data[ServiceAccountJsonKeyFile]; ok {
		return "Google"
	}

	_, first := s.Data[aws.AwsAccessKey]
	_, second := s.Data[aws.AwsSecretKey]
	if first && second {
		return "AWS"
	}

	_, first = s.Data[SSLCertificateChainKey]
	_, second = s.Data[SSLPrivateKeyKey]
	if first && second {
		return "Certificate"
	}
	return "Unknown"
}
