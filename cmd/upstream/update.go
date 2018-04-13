package upstream

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
	"github.com/solo-io/glooctl/pkg/upstream"
	"github.com/solo-io/glooctl/pkg/util"
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
			u, err := runUpdate(sc, filename)
			if err != nil {
				fmt.Printf("Unable to create upstream %q\n", err)
				return
			}
			fmt.Println("Upstream updated")
			output, _ := c.InheritedFlags().GetString("output")
			util.Print(output, "", u,
				func(data interface{}, w io.Writer) error {
					upstream.PrintTable([]*v1.Upstream{data.(*v1.Upstream)}, w)
					return nil
				}, os.Stdout)
		},
	}
	cmd.Flags().StringVarP(&filename, "filename", "f", "", "file to use to create upstream")
	cmd.MarkFlagFilename("filename", "yaml", "yml")
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
	if upstream.Metadata == nil {
		upstream.Metadata = &v1.Metadata{}
	}
	upstream.Metadata.ResourceVersion = existing.Metadata.GetResourceVersion()
	return sc.V1().Upstreams().Update(upstream)
}
