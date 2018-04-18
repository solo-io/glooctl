package secret

import (
	"io/ioutil"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/pkg/errors"
)

const (
	sslCertificateChainKey = "ca_chain"
	sslPrivateKeyKey       = "private_key"
)

type CertificateOptions struct {
	Name       string
	CAChain    string
	PrivateKey string
}

func CreateCertificate(si dependencies.SecretStorage, opts *CertificateOptions) error {
	ca, err := ioutil.ReadFile(opts.CAChain)
	if err != nil {
		return errors.Wrap(err, "unable to read CA chain certificate")
	}
	pk, err := ioutil.ReadFile(opts.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "unable to read private key")
	}
	s := &dependencies.Secret{
		Ref: opts.Name,
		Data: map[string]string{
			sslCertificateChainKey: string(ca),
			sslPrivateKeyKey:       string(pk),
		},
	}
	_, err = si.Create(s)
	return err
}
