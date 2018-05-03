package virtualservice

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/gloo/pkg/storage/file"
	"github.com/solo-io/glooctl/pkg/util"
)

const (
	// DefaultVirtualService is name of the virtual service to create if none exists
	DefaultVirtualService = "default"
)

// Options represents the CLI parameters for virtual services
type Options struct {
	Filename    string
	Output      string
	Template    string
	Interactive bool
}

// ParseFile parses YAML file into a virtual service
func ParseFile(filename string) (*v1.VirtualService, error) {
	var v v1.VirtualService

	// special case: reading from stdin
	if filename == "-" {
		return &v, util.ReadStdinInto(&v)
	}

	err := file.ReadFileInto(filename, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

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

func DefaultVirtualServiceValidation(sc storage.Interface, v *v1.VirtualService) error {
	defaultVS, err := defaultVirtualService(sc, false)
	if err != nil {
		if IsNotExists(err) {
			return nil
		}
		return err
	}

	if v.Name != defaultVS.Name && hasDefaultDomain(v) {
		return fmt.Errorf("not default virtual service needs to specify one or more domains")
	}
	if v.Name == defaultVS.Name && !hasDefaultDomain(v) {
		return fmt.Errorf("default virtual service should have * as one of the domains or be empty")
	}
	return nil
}

func hasDefaultDomain(v *v1.VirtualService) bool {
	if len(v.Domains) == 0 {
		return true
	}

	for _, d := range v.Domains {
		if "*" == d {
			return true
		}
	}
	return false
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
		return nil, NewNotExistsErr("did not find a default virtual service")
	}
	fmt.Println("Did not find a default virtual service. Creating...")
	vservice := &v1.VirtualService{
		Name: DefaultVirtualService,
	}
	return sc.V1().VirtualServices().Create(vservice)
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
