package util

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/protoutil"
)

func ReadStdinInto(v proto.Message) error {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return errors.Errorf("error reading stdin: %v", err)
	}
	jsn, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}
	return protoutil.Unmarshal(jsn, v)
}
