package config

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/solo-io/gloo/pkg/bootstrap"
)

// test to make sure refactoring doesn't switch YAML library
// and generate different YAML format
func TestConfigYAMLFormat(t *testing.T) {
	fh, err := ioutil.TempFile("", "glooctl-config-test")
	if err != nil {
		t.Errorf("unable to create temp file %q", err)
	}
	opts := &bootstrap.Options{}
	defaultConfig(opts)
	err = save(opts, fh.Name())
	if err != nil {
		t.Errorf("error saving config file: %q", err)
	}

	// verify file
	content, err := ioutil.ReadFile(fh.Name())
	if err != nil {
		t.Errorf("unable to read file to verify: %q", err)
	}

	for _, s := range []string{"ConfigStorageOptions:", "Type: kube", "KubeOptions:"} {
		if !bytes.Contains(content, []byte(s)) {
			t.Errorf("doesn't match expected YAML format; didn't find %s", s)
		}
	}
}
