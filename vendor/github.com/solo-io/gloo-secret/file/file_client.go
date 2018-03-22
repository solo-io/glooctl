package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
)

func NewClient(dir string) (secret.SecretInterface, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to setup client with directory %s", dir)
	}

	if !info.IsDir() {
		return nil, errors.New(dir + " is not a directory")
	}
	v := V1{client: &v1Client{dir: dir}}
	return &v, nil
}

type V1 struct {
	client *v1Client
}

func (v *V1) V1() secret.V1 {
	return v.client
}

type v1Client struct {
	dir string
}

func (c *v1Client) Create(s *secret.Secret) (*secret.Secret, error) {
	secrets, err := load(c.dir)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load existing secrets")
	}
	_, exists := secrets[s.Name]
	if exists {
		return nil, errors.Errorf("secret with name %s already exists", s.Name)
	}
	p, d := SecretToFS(s)
	filename := filepath.Join(c.dir, p)
	if err := WriteToFile(filename, d); err != nil {
		return nil, errors.Wrap(err, "unable to save secret")
	}
	return s, nil
}

func (c *v1Client) Update(s *secret.Secret) (*secret.Secret, error) {
	secrets, err := load(c.dir)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load existing secrets")
	}
	_, exists := secrets[s.Name]
	if !exists {
		return nil, errors.Errorf("secret %s does not exist", s.Name)
	}
	p, d := SecretToFS(s)
	filename := filepath.Join(c.dir, p)
	if err := WriteToFile(filename, d); err != nil {
		return nil, errors.Wrap(err, "unable to save secret")
	}
	return s, nil
}

func (c *v1Client) Delete(name string) error {
	secrets, err := load(c.dir)
	if err != nil {
		return errors.Wrap(err, "unable to load existing secrets")
	}
	_, exists := secrets[name]
	if !exists {
		return errors.Errorf("secret %s does not exist", name)
	}
	if err := os.Remove(filepath.Join(c.dir, name)); err != nil {
		return err
	}
	return nil
}

func (c *v1Client) Get(name string) (*secret.Secret, error) {
	filename := filepath.Join(c.dir, name)
	data := make(map[string]string)
	err := ReadFileInto(filename, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get secret %s", name)
	}

	return SecretFromFS(name, data), nil
}

func (c *v1Client) List() ([]*secret.Secret, error) {
	secrets, err := load(c.dir)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load existing secrets")
	}
	out := make([]*secret.Secret, len(secrets))
	i := 0
	for _, v := range secrets {
		out[i] = v
		i++
	}
	return out, nil
}

func load(dir string) (map[string]*secret.Secret, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read directory %s", dir)
	}
	secrets := make(map[string]*secret.Secret)
	for _, f := range files {
		name := f.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		data := make(map[string]string)
		err := ReadFileInto(filepath.Join(dir, name), &data)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse %s as secret", name)
		}
		s := SecretFromFS(name, data)
		secrets[s.Name] = s
	}
	return secrets, nil
}
