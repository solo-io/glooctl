package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/glooctl/pkg/util"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultStorage = "kube"
	configFile     = "config.yaml"
)

// LoadConfig loads saved configuration if any
// if not sets default configuration and also saves it
func LoadConfig(opts *bootstrap.Options) {
	configDir, err := util.ConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to get config directory:", err)
		defaultConfig(opts)
		return
	}
	configFile := filepath.Join(configDir, configFile)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		defaultConfig(opts)
		if os.IsNotExist(err) {
			err = save(opts, configFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unable to save configuration file", configFile)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Error reading configuration file:", err)
		}
		return
	}
	if err := yaml.Unmarshal(data, opts); err != nil {
		defaultConfig(opts)
		fmt.Fprintln(os.Stderr, "Unable to parse configuration file:", err)
	}
}

func SaveConfig(opts *bootstrap.Options) error {
	configDir, err := util.ConfigDir()
	if err != nil {
		errors.Wrap(err, "unable to get glooctl configuration directory")
	}
	err = save(opts, filepath.Join(configDir, configFile))
	if err != nil {
		return errors.Wrap(err, "unable to save configuration")
	}
	return nil
}

func save(opts *bootstrap.Options, configFile string) error {
	b, err := yaml.Marshal(opts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func defaultConfig(opts *bootstrap.Options) {
	opts.ConfigStorageOptions.Type = defaultStorage
	opts.SecretStorageOptions.Type = defaultStorage
	opts.FileStorageOptions.Type = defaultStorage
	opts.KubeOptions.KubeConfig = filepath.Join(util.HomeDir(), ".kube", "config")
}
