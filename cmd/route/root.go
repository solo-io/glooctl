package route

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// RouteCmd returns command related to managing routes on a virtual host
func RouteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "manage routes on a virtual host",
	}

	pflags := cmd.PersistentFlags()
	setupRouteParams(pflags)
	var output string
	pflags.StringVarP(&output, "output", "o", "", "output format yaml|json")
	var domain string
	pflags.StringVarP(&domain, flagDomain, "d", "", "domain for virtual host; empty defaults to default virtual host")
	var vhost string
	pflags.StringVarP(&vhost, flagVirtualHost, "v", "", "specify virtual host by name; empty defaults to default virtual host")
	var file string
	pflags.StringVarP(&file, flagFilename, "f", "", "file with route defintion")
	cmd.MarkFlagFilename(flagFilename)
	cmd.AddCommand(getCmd(), createCmd(), deleteCmd(), sortCmd())
	return cmd
}

func setupRouteParams(flags *pflag.FlagSet) {
	r := routeDetail{}
	flags.StringVarP(&r.event, flagEvent, "e", "", "event type to match")
	flags.StringVar(&r.pathExact, flagPathExact, "", "exact path to match")
	flags.StringVar(&r.pathRegex, flagPathRegex, "", "path regex to match")
	flags.StringVar(&r.pathPrefix, flagPathPrefix, "", "path prefix to match")
	flags.StringVar(&r.verb, flagMethod, "", "HTTP method to match")
	flags.StringVar(&r.headers, flagHeaders, "", "header to match")
	flags.StringVar(&r.upstream, flagUpstream, "", "desitnation upstream")
	flags.StringVar(&r.function, flagFunction, "", "destination function")
}
