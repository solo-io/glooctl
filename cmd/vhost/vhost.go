package vhost

import (
	"fmt"
	"io"
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage/file"
	"github.com/solo-io/gloo/pkg/protoutil"
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

func printJSON(v *v1.VirtualHost) {
	b, err := protoutil.Marshal(v)
	if err != nil {
		fmt.Println("unable to convert to JSON ", err)
		return
	}
	fmt.Println(string(b))
}

func printYAML(v *v1.VirtualHost) {
	writeYAML(v, os.Stdout)
}

func writeYAML(v *v1.VirtualHost, w io.Writer) error {
	jsn, err := protoutil.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "unable to marshal ")
	}
	b, err := yaml.JSONToYAML(jsn)
	if err != nil {
		return errors.Wrap(err, "unable to convert to YAML")
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}

func printJSONList(vhosts []*v1.VirtualHost) {
	for _, v := range vhosts {
		printJSON(v)
	}
}

func printYAMLList(vhosts []*v1.VirtualHost) {
	for _, v := range vhosts {
		printYAML(v)
	}
}

func printSummaryList(v []*v1.VirtualHost) {
	for _, entry := range v {
		fmt.Println(entry.Name)
	}
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
