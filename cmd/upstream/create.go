package upstream

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create upstreams",
		Run: func(c *cobra.Command, args []string) {
			sc, err := util.GetStorageClient(c)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			upstream, err := runCreate(sc, filename)
			if err != nil {
				fmt.Printf("Unable to create upstream %q\n", err)
				return
			}
			fmt.Println("Upstream created")
			output, _ := c.InheritedFlags().GetString("output")
			if output == "yaml" {
				printYAML(upstream)
			}
			if output == "json" {
				printJSON(upstream)
			}
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "file to use to create upstream")
	cmd.MarkFlagFilename("filename")
	cmd.MarkFlagRequired("filename")
	return cmd
}

func runCreate(sc storage.Interface, filename string) (*v1.Upstream, error) {
	upstream, err := parseFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load Upstream from %s", filename)
	}
	return sc.V1().Upstreams().Create(upstream)
}
