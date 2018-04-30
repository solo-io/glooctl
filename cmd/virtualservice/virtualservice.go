package virtualservice

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage/file"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/solo-io/glooctl/pkg/virtualservice"
)

func parseFile(filename string) (*v1.VirtualService, error) {
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

func defaultVirtualServiceValidation(v *v1.VirtualService) error {
	if v.Name != virtualservice.DefaultVirtualService && !hasDomains(v) {
		return fmt.Errorf("not default virtual service needs to specify one or more domains")
	}
	if v.Name == virtualservice.DefaultVirtualService && hasDomains(v) {
		return fmt.Errorf("default virtual service should not have any specific domain")
	}
	return nil
}

func hasDomains(v *v1.VirtualService) bool {
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
