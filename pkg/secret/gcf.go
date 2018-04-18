package secret

import (
	"io/ioutil"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/pkg/errors"
)

const (
	ServiceAccountJsonKeyFile = "json_key_file"
)

type GoogleOptions struct {
	Name     string
	Filename string
}

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
