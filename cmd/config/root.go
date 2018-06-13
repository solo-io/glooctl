package virtualservice

import (
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var (
	filename          string
	overwriteExisting bool
	deleteExisting    bool
)

// Cmd command to manage virtual services
func Cmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "configure",
		Short: `apply an entire configuration (upstreams and virtualservices) to Gloo.

To disable overwriting existing resources defined in the target file, run with -w=false 
To delete any existing resources not defined in the target file, run with -d

The configuration should be provided in a YAML file with the following structure:

upstreams:
- name: foo
  type: service
  spec: ...
- name: bar
  type: kubernetes
  spec: ...
virtual_services:
- name: baz
   domains: ...
- name: qux
   domains: ...

`,
		Run: func(c *cobra.Command, args []string) {
			if err := configure(opts, filename, overwriteExisting, deleteExisting); err != nil {
				fmt.Printf("failed applying configuration %q\n", err)
				os.Exit(1)
			}
		},
	}
	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&filename, "filename", "f", "", "filename to create resources from")
	pflags.BoolVarP(&overwriteExisting, "overwrite", "w", true, "overwrite existing resources "+
		"whose names overlap with those defined in the config file")
	pflags.BoolVarP(&deleteExisting, "delete", "d", false, "delete existing resources "+
		"whose names overlap with those defined in the config file")
	return cmd
}
