package secret

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/plugins/aws"
	"github.com/solo-io/gloo/pkg/storage/dependencies"
)

// AWSOptions represents parameters for creating AWS secret
type AWSOptions struct {
	Name      string
	Filename  string
	KeyID     string
	SecretKey string
	UseEnv    bool
}

// CreateAWS creates secrets that can be used by AWS upstream type
func CreateAWS(si dependencies.SecretStorage, opts *AWSOptions) error {
	id, key, err := idAndKey(opts)
	if err != nil {
		return errors.Wrap(err, "unable to get AWS credentials")
	}
	s := &dependencies.Secret{
		Ref: opts.Name,
		Data: map[string]string{
			aws.AwsAccessKey: id,
			aws.AwsSecretKey: key,
		},
	}
	_, err = si.Create(s)
	return err
}

func idAndKey(opts *AWSOptions) (string, string, error) {
	if opts.KeyID != "" || opts.SecretKey != "" {
		if opts.KeyID != "" && opts.SecretKey != "" {
			return opts.KeyID, opts.SecretKey, nil
		}
		return "", "", errors.New("both access-key-id and secret-access-key must be provided")
	}

	var creds *credentials.Credentials
	if opts.UseEnv {
		creds = credentials.NewEnvCredentials()
	} else {
		//TODO: add a flag for profile
		creds = credentials.NewSharedCredentials(opts.Filename, "")
	}
	vals, err := creds.Get()
	if err != nil {
		return "", "", errors.Wrap(err, "failed to retrieve AWS credentials")
	}
	return vals.AccessKeyID, vals.SecretAccessKey, nil
}
