package vault

import (
	"path/filepath"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
)

func NewClient(vaultAddr, token, basePath string, retries int) (secret.SecretInterface, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = vaultAddr
	cfg.MaxRetries = retries
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "starting vault client")
	}
	client.SetToken(token)

	return &V1{client: &v1Client{client: client, basePath: basePath}}, nil
}

type V1 struct {
	client *v1Client
}

func (v *V1) V1() secret.V1 {
	return v.client
}

type v1Client struct {
	basePath string
	client   *vaultapi.Client
}

func (c *v1Client) Create(s *secret.Secret) (*secret.Secret, error) {
	path, data := SecretToVault(s)
	_, err := c.client.Logical().Write(path, data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create secret")
	}
	return s, nil
}

func (c *v1Client) Update(s *secret.Secret) (*secret.Secret, error) {
	path, data := SecretToVault(s)
	_, err := c.client.Logical().Write(path, data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create secret")
	}
	return s, nil
}

func (c *v1Client) Delete(name string) error {
	_, err := c.client.Logical().Delete(name)
	return err
}

func (c *v1Client) Get(name string) (*secret.Secret, error) {
	s, err := c.client.Logical().Read(name)
	if err != nil {
		return nil, err
	}
	return SecretFromVault(name, s.Data), nil
}

func (c *v1Client) List() ([]*secret.Secret, error) {
	s, err := c.client.Logical().List(c.basePath)
	if err != nil {
		return nil, err
	}
	secrets := make([]*secret.Secret, len(s.Data))
	i := 0
	for k, v := range s.Data {
		data, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.Errorf("unable to convert data for %s", k)
		}
		secrets[i] = SecretFromVault(filepath.Join(c.basePath, k), data)
	}
	return secrets, nil
}
