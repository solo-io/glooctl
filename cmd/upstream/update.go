package upstream

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
	"github.com/spf13/cobra"
)

func updateCmd(opts *client.StorageOptions) *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update upstreams",
		Run: func(c *cobra.Command, args []string) {
			sc, err := client.StorageClient(opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			upstream, err := runUpdate(sc, filename)
			if err != nil {
				fmt.Printf("Unable to create upstream %q\n", err)
				return
			}
			fmt.Println("Upstream updated")
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

func runUpdate(sc storage.Interface, filename string) (*v1.Upstream, error) {
	upstream, err := parseFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load Upstream from %s", filename)
	}
	existing, err := sc.V1().Upstreams().Get(upstream.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find existing Upstream %s", upstream.Name)
	}
	// need to copy new into existing so that we get resource revision number
	copy(upstream, existing)
	return sc.V1().Upstreams().Update(existing)
}

func copy(src *v1.Upstream, dst *v1.Upstream) {
	dst.Type = src.Type
	dst.Spec = src.Spec
	dst.Functions = src.Functions
}
