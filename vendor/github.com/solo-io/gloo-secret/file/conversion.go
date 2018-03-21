package file

import (
	"encoding/base64"
	"log"

	secret "github.com/solo-io/gloo-secret"
)

type fileSecret struct {
	Name            string            `json:"name"`
	Data            map[string]string `json:"data"`
	ResourceVersion string            `json:"version"`
}

func SecretToFS(s *secret.Secret) *fileSecret {
	return &fileSecret{
		Name:            s.Name,
		ResourceVersion: s.ResourceVersion,
		Data:            encode(s.Data),
	}
}

func SecretFromFS(fs *fileSecret) *secret.Secret {
	return &secret.Secret{
		Name:            fs.Name,
		ResourceVersion: fs.ResourceVersion,
		Data:            decode(fs.Data),
	}
}

func encode(in map[string][]byte) map[string]string {
	out := make(map[string]string)
	for k, v := range in {
		out[k] = base64.StdEncoding.EncodeToString(v)
	}
	return out
}

func decode(in map[string]string) map[string][]byte {
	out := make(map[string][]byte)
	for k, v := range in {
		var err error
		out[k], err = base64.StdEncoding.DecodeString(v)
		if err != nil {
			log.Printf("warning unable to decode %s: %q\n", v, err)
		}
	}
	return out
}
