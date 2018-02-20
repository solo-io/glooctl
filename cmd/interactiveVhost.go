package cmd

import (
	"fmt"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
)

func InteractiveModeVhost(cmd string, vh *VHost) {
	if !interactive {
		//return
		fmt.Println("Currently VirtualHost can only be configured in the interactive mode")
	}
	switch cmd {
	case "create":
		fallthrough
	case "update":
		n, _ := getNameAndNamespace(true)
		vh.vhost.Name = n
		getDomains(vh)
		getRoutes(vh)
		getSslConfig(vh)
	case "delete":
		getNameAndNamespace(true)
	case "get":
		fallthrough
	case "describe":
		getNameAndNamespace(false)
	default:
	}
}

func getDomains(vh *VHost) {
	vh.vhost.Domains = make([]string, 0)
	for i := 0; ; i++ {
		s := getString("Domain (\"\" if done)", "", false)
		if *s == "" {
			if i > 0 {
				break
			} else {
				fmt.Println("At least one domain is required")
			}
		} else {
			vh.vhost.Domains = append(vh.vhost.Domains, *s)
		}
	}
}

func getRoutes(vh *VHost) {
	for i := 0; ; i++ {
		cont := getString(fmt.Sprintf("Configuring Route %d. Continue?", i), "yes", true)
		if *cont != "yes" {
			break
		}
		r := &v1.Route{}
		r.Matcher = getMatcher(vh)
		answ := getString("Would you like to configure single destination (1) or multiple weighted destinations (2)", "1", true)
		if *answ == "1" {
			r.SingleDestination = getSingleDestination(vh)
		} else {
			r.MultipleDestinations = getWeightedDestinations(vh)
		}
		r.PrefixRewrite = *getString("Prefix Rewrite", "", false)
	}
}

func getSslConfig(vh *VHost) {
	sref := getString("SSL Secret Reference", "", true)
	vh.vhost.SslConfig = &v1.SSLConfig{SecretRef: *sref}
}

func getMatcher(vh *VHost) *v1.Matcher {
	p := ""
	matcher := v1.Matcher{}
	for {
		answ := *getString("Select matcher path type - Prefix (1), Regex (2), Exact (3)", "1", true)
		switch answ {
		case "1":
			p = *getString("Provide Path Prefix, \"\" to start over", "", false)
			if p == "" {
				continue
			}
			matcher.Path = &v1.Matcher_PathPrefix{PathPrefix: p}
		case "2":
			p = *getString("Provide Path RegEx, \"\" to start over", "", false)
			if p == "" {
				continue
			}
			matcher.Path = &v1.Matcher_PathRegex{PathRegex: p}
		default:
			p = *getString("Provide Exact Path, \"\" to start over", "", false)
			if p == "" {
				continue
			}
			matcher.Path = &v1.Matcher_PathExact{PathExact: p}
		}
		break
	}

	matcher.Headers = make(map[string]string)
	for i := 0; ; i++ {
		k := *getString(fmt.Sprintf("Header %d, \"\" to stop", i), "", false)
		if k == "" {
			break
		}
		v := *getString(fmt.Sprintf("Header %d value", i), "", false)
		matcher.Headers[k] = v
	}

	matcher.QueryParams = make(map[string]string)
	for i := 0; ; i++ {
		k := *getString(fmt.Sprintf("Query Parameter %d, \"\" to stop", i), "", false)
		if k == "" {
			break
		}
		v := *getString(fmt.Sprintf("Query Parameter %d value", i), "", false)
		matcher.QueryParams[k] = v
	}

	matcher.Verbs = make([]string, 0)
	for i := 0; ; i++ {
		v := *getString(fmt.Sprintf("Query Parameter %d, \"\" to stop", i), "", false)
		if v == "" {
			break
		}
		matcher.Verbs = append(matcher.Verbs, v)
	}

	return &matcher
}

func getWeightedDestinations(vh *VHost) []*v1.WeightedDestination {

}

func getSingleDestination(vh *VHost) *v1.Destination {
	des := &v1.Destination{}
	answ := *getString("Select destination type - Function (1) or Upstream (2)", "1", true)

	return des
}
