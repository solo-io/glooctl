package file

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

func WriteToFile(filename string, s *fileSecret) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func ReadFileInto(filename string, s *fileSecret) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Errorf("error reading file: %v", err)
	}
	return json.Unmarshal(data, s)
}
