package cmd

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"
	"github.com/solo-io/glooctl/pkg/function"
	"github.com/spf13/cobra"
)

var (
	functionOpts = struct {
		output   string
		template string
	}{}
)

func functionCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "function",
		Short: "manage functions",
	}
	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&functionOpts.output, "output", "o", "", "output format yaml|json|template")
	pflags.StringVarP(&functionOpts.template, "template", "t", "", "output template")
	cmd.AddCommand(getFunctionsCmd(opts))

	return cmd
}

func getFunctionsCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"list"},
		Short:   "get functions",
		Run: func(c *cobra.Command, a []string) {
			sc, err := configstorage.Bootstrap(*opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			if functionOpts.output == "template" && functionOpts.template == "" {
				fmt.Println("Must provide template when setting output as template")
				return
			}
			if err := function.Get(sc, functionOpts.output, functionOpts.template); err != nil {
				fmt.Printf("Unable to get functions %q\n", err)
				return
			}
		},
	}
	return cmd
}
