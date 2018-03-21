package crd

import (
	secret "github.com/solo-io/gloo-secret"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SecretToCRD(s *secret.Secret) *apiv1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            s.Name,
			ResourceVersion: s.ResourceVersion,
		},
		Data: s.Data,
	}
}

func SecretFromCRD(c *apiv1.Secret) *secret.Secret {
	return &secret.Secret{
		Name:            c.ObjectMeta.Name,
		Data:            c.Data,
		ResourceVersion: c.ResourceVersion,
	}
}
