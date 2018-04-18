package route

import (
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	routeOpt = &routeOption{route: &routeDetail{kube: &kubeUpstream{}}}
)

// RouteCmd returns command related to managing routes on a virtual host
func RouteCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "manage routes on a virtual host",
	}

	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&routeOpt.output, "output", "o", "", "output format yaml|json")
	pflags.StringVarP(&routeOpt.domain, flagDomain, "d", "", "domain for virtual host; empty defaults to default virtual host")
	pflags.StringVarP(&routeOpt.virtualhost, flagVirtualHost, "v", "", "specify virtual host by name; empty defaults to default virtual host")
	pflags.StringVarP(&routeOpt.filename, flagFilename, "f", "", "file with route defintion")
	cmd.MarkFlagFilename(flagFilename, "yaml", "yml")

	create := createCmd(opts)
	update := updateCmd(opts)
	delete := deleteCmd(opts)
	setupRouteParams(create, update, delete)
	cmd.AddCommand(getCmd(opts), create, delete, update, sortCmd(opts))

	annotate(cmd.Flag(flagVirtualHost), "__glooctl_get_virtualhosts")
	return cmd
}

func setupRouteParams(cmds ...*cobra.Command) {
	r := routeOpt.route
	for _, c := range cmds {
		flags := c.Flags()
		flags.StringVarP(&r.event, flagEvent, "e", "", "event type to match")
		flags.StringVar(&r.pathExact, flagPathExact, "", "exact path to match")
		flags.StringVar(&r.pathRegex, flagPathRegex, "", "path regex to match")
		flags.StringVar(&r.pathPrefix, flagPathPrefix, "", "path prefix to match")
		flags.StringVar(&r.verb, flagMethod, "", "HTTP method to match")
		flags.StringVar(&r.headers, flagHeaders, "", "header to match")
		flags.StringVar(&r.upstream, flagUpstream, "", "desitnation upstream")
		flags.StringVar(&r.function, flagFunction, "", "destination function")
		flags.StringVar(&r.prefixRewrite, flagPrefixRewrite, "", "if specified, rewrite the matched portion of "+
			"the path to this value")
		flags.StringVar(&r.extensions, flagExtension, "", "yaml file with route extensions")
		c.MarkFlagFilename(flagExtension, "yaml", "yml")

		// auto complete
		annotate(c.Flag(flagMethod), "__glooctl_route_http_methods")
		annotate(c.Flag(flagUpstream), "__glooctl_get_upstreams")
		annotate(c.Flag(flagFunction), "__glooctl_get_functions")
	}
}

func annotate(f *pflag.Flag, completion string) {
	if f.Annotations == nil {
		f.Annotations = map[string][]string{}
	}
	f.Annotations[cobra.BashCompCustom] = append(f.Annotations[cobra.BashCompCustom], completion)
}
