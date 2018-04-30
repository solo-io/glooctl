package virtualservice

import (
	"fmt"
	"strings"

	"github.com/solo-io/gloo/pkg/storage/dependencies"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
	psecret "github.com/solo-io/glooctl/pkg/secret"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// SelectInteractive - selects a virtual service interactively from the list
// of virtual services in the system
func SelectInteractive(sc storage.Interface) (*v1.VirtualService, error) {
	existing, err := sc.V1().VirtualServices().List()
	if err != nil {
		return nil, err
	}
	if len(existing) == 0 {
		return nil, errors.New("no existing virtual services to select")
	}

	virtualServiceNames := make([]string, len(existing))
	for i, v := range existing {
		virtualServiceNames[i] = v.Name
	}

	var selected string
	if err := survey.AskOne(&survey.Select{
		Message: "Please select the virtual service to edit:",
		Options: virtualServiceNames,
	}, &selected, survey.Required); err != nil {
		return nil, err
	}

	for _, v := range existing {
		if v.Name == selected {
			return v, nil
		}
	}
	return nil, errors.New("didn't find selected virtual service")
}

// Interactive for interactively creating/editing virtual services
// Doesn't handle routes as we have separate interactive mode for routes
func Interactive(sc storage.Interface, si dependencies.SecretStorage, vs *v1.VirtualService) error {
	existing, err := sc.V1().VirtualServices().List()
	if err != nil {
		return err
	}
	// name
	if vs.Name == "" {
		// new virtual service
		var name string
		err = survey.AskOne(&survey.Input{
			Message: "Please enter a name for virtual service:",
		}, &name, func(val interface{}) error {
			v, ok := val.(string)
			if !ok {
				return errors.New("not a string value")
			}
			if v == "" {
				return errors.New("virtual service name can't be empty")
			}
			for _, e := range existing {
				if e.Name == v {
					return errors.New("virtual service with that name already exists")
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		vs.Name = name
	}

	updatedDomains, err := domainsInteractive(vs.Domains)
	if err != nil {
		return err
	}
	vs.Domains = updatedDomains

	updatedSSL, err := sslConfigInteractive(si, vs.SslConfig)
	if err != nil {
		return err
	}
	vs.SslConfig = updatedSSL
	return nil
}

func domainsInteractive(list []string) ([]string, error) {
	if len(list) != 0 {
		printDomains(list)
		replace := false
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to replace the existing domains?"},
			&replace, nil); err != nil {
			return nil, err
		}
		if !replace {
			return list, nil
		}
	} else {
		set := false
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to set domains?"},
			&set, nil); err != nil {
			return nil, err
		}
		if !set {
			return nil, nil
		}
	}

	var newDomains string
	if err := survey.AskOne(&survey.Input{Message: "Please enter a comma separated list of domains (leave empty to set none):"},
		&newDomains, nil); err != nil {
		return nil, err
	}
	if newDomains == "" {
		return nil, nil
	}

	return strings.Split(newDomains, ","), nil
}

func sslConfigInteractive(si dependencies.SecretStorage, ssl *v1.SSLConfig) (*v1.SSLConfig, error) {
	if ssl != nil {
		printSSLConfig(ssl)
		replace := false
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to replace the existing SSL configuration?"},
			&replace, nil); err != nil {
			return nil, err
		}
		if !replace {
			return ssl, nil
		}
	} else {
		set := false
		if err := survey.AskOne(&survey.Confirm{Message: "Do you want to set SSL configuration?"},
			&set, nil); err != nil {
			return nil, err
		}
		if !set {
			return nil, nil
		}
	}

	secrets, err := psecret.SecretRefs(si, isCertificate)
	if err != nil {
		return nil, err
	}
	if len(secrets) == 0 {
		return nil, errors.New("unable to get secret reference for certificates")
	}

	var secretRef string
	secretOpts := append([]string{"None"}, secrets...)
	if err := survey.AskOne(&survey.Select{
		Message: "Please select a secret reference for certificate",
		Options: secretOpts,
	}, &secretRef, survey.Required); err != nil {
		return nil, err
	}

	if "None" == secretRef {
		return nil, nil
	}
	return &v1.SSLConfig{SecretRef: secretRef}, nil
}

func isCertificate(s *dependencies.Secret) bool {
	if s.Data == nil {
		return false
	}

	_, first := s.Data[psecret.SSLCertificateChainKey]
	_, second := s.Data[psecret.SSLPrivateKeyKey]
	return first && second
}

func printSSLConfig(ssl *v1.SSLConfig) {
	fmt.Printf("Secret Ref for SSL: %s\n\n\n", ssl.SecretRef)
}

func printDomains(list []string) {
	fmt.Printf("Domains: %s\n\n\n", strings.Join(list, ", "))
}
