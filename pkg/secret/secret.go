package secret

import (
	"os"

	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
	"github.com/solo-io/gloo/pkg/storage"
)

func Get(sc storage.Interface, si secret.SecretInterface, name string) error {
	var list []*secret.Secret
	if name != "" {
		s, err := si.V1().Get(name)
		if err != nil {
			return errors.Wrap(err, "unable to get secret "+name)
		}
		list = []*secret.Secret{s}
	} else {
		var err error
		list, err = si.V1().List()
		if err != nil {
			return errors.Wrap(err, "unable to get secrets")
		}
	}
	upstreams, err := sc.V1().Upstreams().List()
	if err != nil {
		return errors.Wrap(err, "unable to get upstreams")
	}
	virtualhosts, err := sc.V1().VirtualHosts().List()
	if err != nil {
		return errors.Wrap(err, "unable to get virtual hosts")
	}

	PrintTableWithUsage(list, os.Stdout, upstreams, virtualhosts)
	return nil
}
