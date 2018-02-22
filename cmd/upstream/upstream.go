package upstream

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage/file"
	"github.com/solo-io/gloo/pkg/protoutil"
)

func parseFile(filename string) (*v1.Upstream, error) {
	var u v1.Upstream
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
	jsn, err := protoutil.Marshal(u)
	if err != nil {
		fmt.Println("unable to marshal ", err)
		return
	}
	b, err := yaml.JSONToYAML(jsn)
	if err != nil {
		fmt.Println("unable to convert to YAML ", err)
		return
	}
	fmt.Println(string(b))
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
