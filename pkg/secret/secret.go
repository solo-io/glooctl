package secret

import (
	"os"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/storage"
)

// Get gets a secret with given name or all secrets if no name is provided and prints
// them and their usage
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
	virtualservices, err := sc.V1().VirtualServices().List()
	if err != nil {
		return errors.Wrap(err, "unable to get virtual services")
	}

	PrintTableWithUsage(list, os.Stdout, upstreams, virtualservices)
	return nil
}

// SecretRefs returns a list of secret references filtered using the provided filter
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
