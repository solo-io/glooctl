package file

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

func WriteToFile(filename string, data map[string]string) error {
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func ReadFileInto(filename string, data *map[string]string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Errorf("error reading file: %v", err)
	}
	return yaml.Unmarshal(b, data)
}
