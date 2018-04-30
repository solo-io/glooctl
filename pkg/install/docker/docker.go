package docker

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/glooctl/pkg/config"

	"github.com/pkg/errors"
)

func Install(folder string) error {
	err := createInstallFolder(folder)
	if err != nil {
		return err
	}
	err = writeDockerCompose(folder)
	if err != nil {
		return err
	}
	err = writeEnvoyConfig(folder)
	if err != nil {
		return err
	}

	storageFolder, err := createStorageFolders(folder)
	if err != nil {
		return err
	}

	return updateGlooctlConfig(storageFolder)
}

func createInstallFolder(folder string) error {
	stat, err := os.Stat(folder)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "unable to setup install directory")
		}
		err = os.MkdirAll(folder, 0755)
		if err != nil {
			return errors.Wrap(err, "unable to create directory")
		}
		return nil
	}

	if !stat.IsDir() {
		return errors.Errorf("%s already exists and isn't a directory", folder)
	}

	return nil
}

func writeDockerCompose(folder string) error {
	filename := filepath.Join(folder, "docker-compose.yml")
	err := ioutil.WriteFile(filename, []byte(dockerComposeYAML), 0644)
	if err != nil {
		return errors.Wrap(err, "unable to create "+filename)
	}
	return nil
}

func writeEnvoyConfig(folder string) error {
	filename := filepath.Join(folder, "envoy.yaml")
	err := ioutil.WriteFile(filename, []byte(envoyYAML), 0644)
	if err != nil {
		return errors.Wrap(err, "unable to create "+filename)
	}
	return nil
}

func createStorageFolders(folder string) (string, error) {
	glooConfig := filepath.Join(folder, "gloo-config")
	err := os.MkdirAll(glooConfig, 0755)
	if err != nil {
		return "", errors.Wrap(err, "unable to create storage directory "+glooConfig)
	}

	for _, f := range []string{"_gloo_config", "_gloo_secrets", "_gloo_files"} {
		err = os.MkdirAll(filepath.Join(glooConfig, f), 0755)
		if err != nil {
			return "", errors.Wrap(err, "unable to create storage directory "+f)
		}
	}

	return glooConfig, nil
}

func updateGlooctlConfig(folder string) error {
	opts := &bootstrap.Options{}

	opts.ConfigStorageOptions.Type = "file"
	opts.FileStorageOptions.Type = "file"
	opts.SecretStorageOptions.Type = "file"

	opts.FileOptions.ConfigDir = filepath.Join(folder, "_gloo_config")
	opts.FileOptions.FilesDir = filepath.Join(folder, "_gloo_files")
	opts.FileOptions.SecretDir = filepath.Join(folder, "_gloo_secrets")

	err := config.SaveConfig(opts)
	if err != nil {
		return errors.Wrap(err, "unable to configure glooctl")
	}

	return nil
}
