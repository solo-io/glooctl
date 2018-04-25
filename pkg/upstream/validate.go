package upstream

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugins/aws"
	"github.com/solo-io/gloo/pkg/plugins/google"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/gloo/pkg/storage/dependencies"
	psecret "github.com/solo-io/glooctl/pkg/secret"
)

type upstreamValidator func(storage.Interface, dependencies.SecretStorage, *v1.Upstream) (bool, string)
type validatorPlugin struct {
	upstreamType string
	validator    upstreamValidator
}

var (
	validatorPlugins = []validatorPlugin{
		{aws.UpstreamTypeAws, validateAWS},
		{gfunc.UpstreamTypeGoogle, validateGoogle},
	}
)

func Validate(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) (bool, string) {
	for _, p := range validatorPlugins {
		if p.upstreamType == u.Type {
			return p.validator(sc, si, u)
		}
	}
	// unknown upstream type; no validator
	// assume it is valid
	return true, ""
}

func validateAWS(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) (bool, string) {
	lambdaSpec, err := aws.DecodeUpstreamSpec(u.Spec)
	if err != nil {
		return false, fmt.Sprintf("unable to decode AWS upstream: %q", err)
	}
	awsSecrets, err := si.Get(lambdaSpec.SecretRef)
	if err != nil {
		// warning
		return true, fmt.Sprintf("did not verify referenced secret")
	}
	if _, ok := awsSecrets.Data[aws.AwsAccessKey]; !ok {
		return false, fmt.Sprintf("AWS Access Key missing in referenced secret")
	}
	if _, ok := awsSecrets.Data[aws.AwsSecretKey]; !ok {
		return false, fmt.Sprintf("AWS Secret Key missing in referenced secret")
	}
	return true, ""
}

func validateGoogle(sc storage.Interface, si dependencies.SecretStorage, u *v1.Upstream) (bool, string) {
	_, err := gfunc.DecodeUpstreamSpec(u.Spec)
	if err != nil {
		return false, fmt.Sprintf("unable to decode GCF upstream: %q", err)
	}
	secretRef, ok := u.Metadata.Annotations[psecret.GoogleAnnotationKey]
	if !ok {
		return true, fmt.Sprintf("Google Cloud Function Discovery requires annotation wity key %s", psecret.GoogleAnnotationKey)
	}

	gcfSecret, err := si.Get(secretRef)
	if err != nil {
		return true, fmt.Sprintf("unable to verify referenced secret")
	}
	if _, ok := gcfSecret.Data[psecret.ServiceAccountJsonKeyFile]; !ok {
		return false, fmt.Sprintf("secret missing key %s", psecret.ServiceAccountJsonKeyFile)
	}
	return true, ""
}
