package file

import (
	secret "github.com/solo-io/gloo-secret"
)

func SecretFromFS(path string, data map[string]string) *secret.Secret {
	return &secret.Secret{
		Name: path,
		Data: decode(data),
	}
}

func SecretToFS(s *secret.Secret) (string, map[string]string) {
	return s.Name, encode(s.Data)
}

func decode(in map[string]string) map[string][]byte {
	out := make(map[string][]byte)
	for k, v := range in {
		out[k] = []byte(v)
	}
	return out
}

func encode(in map[string][]byte) map[string]string {
	out := make(map[string]string)
	for k, v := range in {
		out[k] = string(v)
	}
	return out
}
