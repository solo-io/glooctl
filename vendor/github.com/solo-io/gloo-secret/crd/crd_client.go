package crd

import (
	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

const GlooDefaultNamespace = "gloo-system"

func NewClient(cfg *rest.Config, namespace string) (secret.SecretInterface, error) {
	if namespace == "" {
		namespace = GlooDefaultNamespace
	}
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get client")
	}
	si := cs.CoreV1().Secrets(namespace)
	return &V1{client: &v1Client{si: si}}, nil
}

type V1 struct {
	client *v1Client
}

func (v *V1) V1() secret.V1 {
	return v.client
}

type v1Client struct {
	si corev1.SecretInterface
}

func (c *v1Client) Create(s *secret.Secret) (*secret.Secret, error) {
	created, err := c.si.Create(SecretToCRD(s))
	if err != nil {
		return nil, err
	}
	return SecretFromCRD(created), nil
}

func (c *v1Client) Update(s *secret.Secret) (*secret.Secret, error) {
	updated, err := c.si.Update(SecretToCRD(s))
	if err != nil {
		return nil, err
	}
	return SecretFromCRD(updated), nil
}

func (c *v1Client) Delete(name string) error {
	return c.si.Delete(name, &metav1.DeleteOptions{})
}

func (c *v1Client) Get(name string) (*secret.Secret, error) {
	crdSecret, err := c.si.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return SecretFromCRD(crdSecret), nil
}

func (c *v1Client) List() ([]*secret.Secret, error) {
	list, err := c.si.List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	secrets := make([]*secret.Secret, len(list.Items))
	for i, s := range list.Items {
		secrets[i] = SecretFromCRD(&s)
	}
	return secrets, nil
}
