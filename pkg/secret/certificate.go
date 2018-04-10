package secret

import (
	"io/ioutil"

	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
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

func CreateCertificate(si secret.SecretInterface, opts *CertificateOptions) error {
	ca, err := ioutil.ReadFile(opts.CAChain)
	if err != nil {
		return errors.Wrap(err, "unable to read CA chain certificate")
	}
	pk, err := ioutil.ReadFile(opts.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "unable to read private key")
	}
	s := &secret.Secret{
		Name: opts.Name,
		Data: map[string][]byte{
			sslCertificateChainKey: ca,
			sslPrivateKeyKey:       pk,
		},
	}
	_, err = si.V1().Create(s)
	return err
}
