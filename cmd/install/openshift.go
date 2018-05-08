package install

import (
	"fmt"
	"os"

	"github.com/solo-io/glooctl/pkg/install/openshift"
	"github.com/spf13/cobra"
)

func openshiftCmd() *cobra.Command {
	dryRun := false
	cmd := &cobra.Command{
		Use:   "openshift",
		Short: "install gloo on OpenShift",
		Long: `
	Installs latest gloo on OpenShift. It downloads the latest installation YAML
	file and installs to the current OpenShift context.`,
		Run: func(c *cobra.Command, a []string) {
			err := openshift.Install(dryRun)
			if err != nil {
				fmt.Printf("Unable to isntall gloo on OpenShift %q\n", err)
				os.Exit(1)
			}
			if !dryRun {
				fmt.Println("Gloo successfully installed.")
			}
		},
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false,
		"If true, only print the objects that will be setup, without sending it")
	return cmd
}
