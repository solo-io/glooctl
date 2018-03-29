package vault

import (
	"log"

	secret "github.com/solo-io/gloo-secret"
)

func SecretFromVault(path string, data map[string]interface{}) *secret.Secret {
	return &secret.Secret{
		Name: path,
		Data: decode(data),
	}
}

func SecretToVault(s *secret.Secret) (string, map[string]interface{}) {
	return s.Name, encode(s.Data)
}

func decode(in map[string]interface{}) map[string][]byte {
	out := make(map[string][]byte)
	for k, v := range in {
		s, ok := v.(string)
		if !ok {
			log.Println("warning: value not a string for key", k)
		}
		out[k] = []byte(s)
	}
	return out
}

func encode(in map[string][]byte) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range in {
		out[k] = string(v)
	}
	return out
}
