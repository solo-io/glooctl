package upstream

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

func parseFile(filename string) (*v1.Upstream, error) {
	var u v1.Upstream
	// special case: reading from stdin
	if filename == "-" {
		return &u, util.ReadStdinInto(&u)
	}
	err := file.ReadFileInto(filename, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func printJSON(u *v1.Upstream) {
	b, err := protoutil.Marshal(u)
	if err != nil {
		fmt.Println("unable to convert to JSON ", err)
		return
	}
	fmt.Println(string(b))
}

func printYAML(u *v1.Upstream) {
	writeYAML(u, os.Stdout)
}

func writeYAML(u *v1.Upstream, w io.Writer) error {
	jsn, err := protoutil.Marshal(u)
	if err != nil {
		return errors.Wrap(err, "unable to marshal")
	}
	b, err := yaml.JSONToYAML(jsn)
	if err != nil {
		return errors.Wrap(err, "unable to convert to YAML")
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}

func printJSONList(u []*v1.Upstream) {
	for _, entry := range u {
		printJSON(entry)
	}
}

func printYAMLList(u []*v1.Upstream) {
	for _, entry := range u {
		printYAML(entry)
	}
}

func printSummaryList(u []*v1.Upstream) {
	for _, entry := range u {
		fmt.Println(entry.Name)
	}
}
