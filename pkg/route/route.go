package route

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
)

const (
	// DefaultVirtualHost is name of the virtual host to create if none exists
	DefaultVirtualHost = "default"
)

// VirtualHost returns a virtual host for given name or domain
// If name or domain isn't provided, it returns default virtual host (if exits or create is true)
func VirtualHost(sc storage.Interface, vhostname, domain string, create bool) (*v1.VirtualHost, error) {
	if vhostname != "" {
		vh, err := sc.V1().VirtualHosts().Get(vhostname)
		if err != nil {
			return nil, err
		}
		return vh, nil
	}

	if domain != "" {
		// find all virtual hosts that can match
		virtualHosts, err := sc.V1().VirtualHosts().List()
		if err != nil {
			return nil, errors.Wrap(err, "unable to get list of virtual hosts")
		}
		virtualHosts = virtualHostsForDomain(virtualHosts, domain)
		switch len(virtualHosts) {
		case 0:
			// TODO? if create is true, should we create a new virtual host with the domain?
			// should we add this domain to default virtual host?
			return nil, fmt.Errorf("didn't find any virtual host for the domain %s", domain)
		case 1:
			return virtualHosts[0], nil
		default:
			return nil, fmt.Errorf("the domain %s matched %d virtual hosts", domain, len(virtualHosts))
		}
	}

	return defaultVirtualHost(sc, create)
}

func contains(vh *v1.VirtualHost, d string) bool {
	for _, e := range vh.Domains {
		if e == d {
			return true
		}
	}
	return false
}

func virtualHostsForDomain(virtualHosts []*v1.VirtualHost, domain string) []*v1.VirtualHost {
	var validOnes []*v1.VirtualHost
	for _, v := range virtualHosts {
		if contains(v, domain) {
			validOnes = append(validOnes, v)
		}
	}
	return validOnes
}

func defaultVirtualHost(sc storage.Interface, create bool) (*v1.VirtualHost, error) {
	// does one exist?
	vhosts, err := sc.V1().VirtualHosts().List()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of existing virtual hosts")
	}
	for _, v := range vhosts {
		if v.Domains == nil ||
			len(v.Domains) == 0 ||
			contains(v, "*") {
			return v, nil
		}
	}

	if !create {
		return nil, fmt.Errorf("did not find a default virtual host")
	}
	fmt.Println("Did not find a default virtual host. Creating...")
	vhost := &v1.VirtualHost{
		Name: DefaultVirtualHost,
	}
	return sc.V1().VirtualHosts().Create(vhost)
}
