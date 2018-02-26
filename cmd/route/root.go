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
	pflags.StringVarP(&domain, "domain", "d", "", "domain for virtual host; empty defaults to default virtual host")
	var file string
	pflags.StringVarP(&file, "filename", "f", "", "file with route defintion")
	cmd.MarkFlagFilename("filename")
	cmd.AddCommand(getCmd(), createCmd(), deleteCmd(), sortCmd())
	return cmd
}

func setupRouteParams(flags *pflag.FlagSet) {
	r := routeDetail{}
	flags.StringVarP(&r.event, "event", "e", "", "event type to match")
	flags.StringVar(&r.pathExact, "path-exact", "", "exact path to match")
	flags.StringVar(&r.pathRegex, "path-regex", "", "path regex to match")
	flags.StringVar(&r.pathPrefix, "path-prefix", "", "path prefix to match")
	flags.StringVar(&r.verb, "http-method", "", "HTTP method to match")
	flags.StringVar(&r.headers, "header", "", "header to match")
	flags.StringVar(&r.upstream, "upstream", "", "desitnation upstream")
	flags.StringVar(&r.function, "function", "", "destination function")
}
