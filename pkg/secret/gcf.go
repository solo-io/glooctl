package secret

import (
	"io/ioutil"

	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
)

const (
	ServiceAccountJsonKeyFile = "json_key_file"
)

type GoogleOptions struct {
	Name     string
	Filename string
}

func CreateGoogle(si secret.SecretInterface, opts *GoogleOptions) error {
	b, err := ioutil.ReadFile(opts.Filename)
	if err != nil {
		return errors.Wrapf(err, "unable to read service account key file %s", opts.Filename)
	}
	s := &secret.Secret{
		Name: opts.Name,
		Data: map[string][]byte{
			ServiceAccountJsonKeyFile: b,
		},
	}
	_, err = si.V1().Create(s)
	return err
}
