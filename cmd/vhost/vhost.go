package vhost

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage/file"
	"github.com/solo-io/glooctl/pkg/util"
)

const defaultVHost = "default"

func parseFile(filename string) (*v1.VirtualHost, error) {
	var v v1.VirtualHost

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

func defaultVHostValidation(v *v1.VirtualHost) error {
	if v.Name != defaultVHost && !hasDomains(v) {
		return fmt.Errorf("not default virtual host needs to specify one or more domains")
	}
	if v.Name == defaultVHost && hasDomains(v) {
		return fmt.Errorf("default virtual host should not have any specific domain")
	}
	return nil
}

func hasDomains(v *v1.VirtualHost) bool {
	if v.Domains == nil {
		return false
	}
	if len(v.Domains) == 0 {
		return false
	}

	if len(v.Domains) == 1 {
		return "*" != v.Domains[0]
	}
	return true
}
