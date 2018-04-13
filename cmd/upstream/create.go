package upstream

// FIXME - replace kube secret interface with secrets client
import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugins/aws"
	"github.com/solo-io/gloo/pkg/plugins/google"
	storage "github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/glooctl/pkg/client"
	"github.com/solo-io/glooctl/pkg/upstream"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

const (
	// expected map identifiers for secrets
	awsAccessKey = "access_key"
	awsSecretKey = "secret_key"

	annotationKey = "gloo.solo.io/google_secret_ref"
	// expected map identifiers for secrets
	serviceAccountJsonKeyFile = "json_key_file"
)

func createCmd(opts *client.StorageOptions) *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create upstreams",
		Run: func(c *cobra.Command, args []string) {
			sc, err := client.StorageClient(opts)
			if err != nil {
				fmt.Printf("Unable to create storage client %q\n", err)
				return
			}
			si, err := client.SecretClient(opts)
			if err != nil {
				fmt.Printf("Unable to create secret client %q\n", err)
				return
			}
			u, err := runCreate(sc, si, filename)
			if err != nil {
				fmt.Printf("Unable to create upstream %q\n", err)
				return
			}
			fmt.Println("Upstream created")
			output, _ := c.InheritedFlags().GetString("output")
			util.Print(output, tplt, u,
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

func runCreate(sc storage.Interface, si secret.SecretInterface, filename string) (*v1.Upstream, error) {
	upstream, err := parseFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load Upstream from %s", filename)
	}

	valid, message := validate(sc, si, upstream)
	if !valid {
		return nil, fmt.Errorf("invalid upstream: %s", message)
	}
	// add verbose mode to disable this normally
	if message != "" {
		fmt.Println("Warning:", message)
	}
	return sc.V1().Upstreams().Create(upstream)
}

func validate(sc storage.Interface, si secret.SecretInterface, u *v1.Upstream) (bool, string) {
	switch u.Type {
	case aws.UpstreamTypeAws:
		lambdaSpec, err := aws.DecodeUpstreamSpec(u.Spec)
		if err != nil {
			return false, fmt.Sprintf("Unable to decode lambda upstream spec: %q", err)
		}
		awsSecrets, err := si.V1().Get(lambdaSpec.SecretRef)
		if err != nil {
			// warning
			return true, fmt.Sprintf("Unable to load referenced secret. Please make sure it exists.")
		}
		if _, ok := awsSecrets.Data[awsAccessKey]; !ok {
			return false, fmt.Sprintf("AWS Access Key missing in referenced secret")
		}
		if _, ok := awsSecrets.Data[awsSecretKey]; !ok {
			return false, fmt.Sprintf("AWS Secret Key missing in referenced secret")
		}
		return true, ""
	case gfunc.UpstreamTypeGoogle:
		_, err := gfunc.DecodeUpstreamSpec(u.Spec)
		if err != nil {
			return false, fmt.Sprintf("Unable to decode GCF upstream spec: %q", err)
		}
		secretRef, ok := u.Metadata.Annotations[annotationKey]
		if !ok {
			return true, fmt.Sprintf("Google Cloud Function Discovery requires annotation wity key %s.", annotationKey)
		}

		gcfSecret, err := si.V1().Get(secretRef)
		if err != nil {
			return true, fmt.Sprintf("Unable to verify referenced secret. Please make sure it exists.")
		}
		if _, ok := gcfSecret.Data[serviceAccountJsonKeyFile]; !ok {
			return false, fmt.Sprintf("secret missing key %s", serviceAccountJsonKeyFile)
		}
		return true, ""
	default:
		return true, ""
	}
}
