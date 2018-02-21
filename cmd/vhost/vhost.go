package vhost

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage/file"
)

func parseFile(filename string) (*v1.VirtualHost, error) {
	var v v1.VirtualHost
	err := file.ReadFileInto(filename, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func printJSON(v *v1.VirtualHost) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("unable to convert to JSON ", err)
		return
	}
	fmt.Println(string(b))
}

func printYAML(v *v1.VirtualHost) {
	b, err := yaml.Marshal(v)
	if err != nil {
		fmt.Println("unable to convert to YAML", err)
		return
	}
	fmt.Println(string(b))
}

func printJSONList(v []*v1.VirtualHost) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("unable to convert to JSON ", err)
		return
	}
	fmt.Println(string(b))
}

func printYAMLList(v []*v1.VirtualHost) {
	b, err := yaml.Marshal(v)
	if err != nil {
		fmt.Println("unable to convert to YAML ", err)
		return
	}
	fmt.Println(string(b))
}

func printSummaryList(v []*v1.VirtualHost) {
	for _, entry := range v {
		fmt.Println(entry.Name)
	}
}
