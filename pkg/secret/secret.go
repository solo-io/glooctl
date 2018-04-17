package secret

import (
	"os"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/storage"
)

func Get(sc storage.Interface, si dependencies.SecretStorage, name string) error {
	var list []*dependencies.Secret
	if name != "" {
		s, err := si.Get(name)
		if err != nil {
			return errors.Wrap(err, "unable to get secret "+name)
		}
		list = []*dependencies.Secret{s}
	} else {
		var err error
		list, err = si.List()
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

func SecretRefs(si dependencies.SecretStorage, filter func(*dependencies.Secret) bool) ([]string, error) {
	secrets, err := si.List()
	if err != nil {
		return nil, err
	}
	var refs []string
	for _, s := range secrets {
		if filter(s) {
			refs = append(refs, s.Ref)
		}
	}
	return refs, nil
}
