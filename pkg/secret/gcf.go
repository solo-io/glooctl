package secret

import (
	"io/ioutil"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/pkg/errors"
)

const (
	ServiceAccountJsonKeyFile = "json_key_file"
)

// GoogleOptions represents the parameters needed to create secret for Google
type GoogleOptions struct {
	Name     string
	Filename string
}

// CreateGoogle creates a secret for function discovery service to use
func CreateGoogle(si dependencies.SecretStorage, opts *GoogleOptions) error {
	b, err := ioutil.ReadFile(opts.Filename)
	if err != nil {
		return errors.Wrapf(err, "unable to read service account key file %s", opts.Filename)
	}
	s := &dependencies.Secret{
		Ref: opts.Name,
		Data: map[string]string{
			ServiceAccountJsonKeyFile: string(b),
		},
	}
	_, err = si.Create(s)
	return err
}
