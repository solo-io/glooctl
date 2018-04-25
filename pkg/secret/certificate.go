package secret

import (
	"io/ioutil"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/pkg/errors"
)

const (
	// SSLCertificateChainKey represents the key to identify certificate chain in the secret
	SSLCertificateChainKey = "ca_chain"
	// SSLPrivateKeyKey represents the key to identify private key in the secret
	SSLPrivateKeyKey = "private_key"
)

// CertificateOptions represents parameters for creating secrets representing SSL certificates
type CertificateOptions struct {
	Name       string
	CAChain    string
	PrivateKey string
}

// CreateCertificate creates a secret representing SSL certificate
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
			SSLCertificateChainKey: string(ca),
			SSLPrivateKeyKey:       string(pk),
		},
	}
	_, err = si.Create(s)
	return err
}
