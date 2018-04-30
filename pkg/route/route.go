package route

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
)

const (
	// DefaultVirtualService is name of the virtual service to create if none exists
	DefaultVirtualService = "default"
)

// VirtualService returns a virtual service for given name or domain
// If name or domain isn't provided, it returns default virtual service (if exits or create is true)
func VirtualService(sc storage.Interface, virtualServiceName, domain string, create bool) (*v1.VirtualService, error) {
	if virtualServiceName != "" {
		vs, err := sc.V1().VirtualServices().Get(virtualServiceName)
		if err != nil {
			return nil, err
		}
		return vs, nil
	}

	if domain != "" {
		// find all virtual services that can match
		virtualServices, err := sc.V1().VirtualServices().List()
		if err != nil {
			return nil, errors.Wrap(err, "unable to get list of virtual services")
		}
		virtualServices = virtualServicesForDomain(virtualServices, domain)
		switch len(virtualServices) {
		case 0:
			return nil, fmt.Errorf("didn't find any virtual service for the domain %s", domain)
		case 1:
			return virtualServices[0], nil
		default:
			return nil, fmt.Errorf("the domain %s matched %d virtual services", domain, len(virtualServices))
		}
	}

	return defaultVirtualService(sc, create)
}

func contains(vs *v1.VirtualService, d string) bool {
	for _, e := range vs.Domains {
		if e == d {
			return true
		}
	}
	return false
}

func virtualServicesForDomain(virtualServices []*v1.VirtualService, domain string) []*v1.VirtualService {
	var validOnes []*v1.VirtualService
	for _, v := range virtualServices {
		if contains(v, domain) {
			validOnes = append(validOnes, v)
		}
	}
	return validOnes
}

func defaultVirtualService(sc storage.Interface, create bool) (*v1.VirtualService, error) {
	// does one exist?
	vservices, err := sc.V1().VirtualServices().List()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of existing virtual services")
	}
	for _, v := range vservices {
		if v.Domains == nil ||
			len(v.Domains) == 0 ||
			contains(v, "*") {
			return v, nil
		}
	}

	if !create {
		return nil, fmt.Errorf("did not find a default virtual service")
	}
	fmt.Println("Did not find a default virtual service. Creating...")
	vservice := &v1.VirtualService{
		Name: DefaultVirtualService,
	}
	return sc.V1().VirtualServices().Create(vservice)
}
