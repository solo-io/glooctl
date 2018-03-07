package route

import (
	"fmt"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	var sort bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update a route",
		Long: `
Update a route based on either the definition in the YAML file
or based on the route matcher and destination provided in the CLI.

While selecting route to update, glooctl matches routes based on
matcher and destination only. It doesn't include extensions.`,
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}

			flags := c.InheritedFlags()
			domain, _ := flags.GetString("domain")
			vhostname, _ := flags.GetString(flagVirtualHost)
			route, err := route(c.Flags(), sc)
			if err != nil {
				fmt.Printf("Unable to get route %q\n", err)
				return
			}
			routes, err := runUpdate(sc, vhostname, domain, route, sort)
			if err != nil {
				fmt.Printf("Unable to get route for %s: %q\n", domain, err)
			}
			output, _ := flags.GetString("output")
			printRoutes(routes, output)
		},
	}
	kube := kubeUpstream{}
	flags := cmd.Flags()
	flags.StringVar(&kube.name, flagKubeName, "", "kubernetes service name")
	flags.StringVar(&kube.namespace, flagKubeNamespace, "", "kubernetes service namespace")
	flags.IntVar(&kube.port, flagKubePort, 0, "kubernetes service port")
	flags.BoolVar(&sort, "sort", false, "sort the routes after appending the new route")
	return cmd
}

func runUpdate(sc storage.Interface, vhostname, domain string, route *v1.Route, sort bool) ([]*v1.Route, error) {
	v, created, err := virtualHost(sc, vhostname, domain, false)
	if err != nil {
		return nil, err
	}
	if created {
		fmt.Println("Using newly created virtual host:", v.Name)
	} else {
		fmt.Println("Using virtual host:", v.Name)
	}

	existing := v.GetRoutes()
	updated := make([]*v1.Route, len(existing))
	var matches []*v1.Route
	for i, r := range existing {
		if match(route, r) {
			matches = append(matches, r)
			updated[i] = route
			continue
		}
		updated[i] = r
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("The given route parameters matches more than one route.")
	}

	v.Routes = updated
	if sort {
		sortRoutes(v.Routes)
	}

	vh, err := sc.V1().VirtualHosts().Update(v)
	if err != nil {
		return nil, err
	}
	return vh.GetRoutes(), nil
}
