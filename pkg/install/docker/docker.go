package docker

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/glooctl/pkg/config"

	"github.com/pkg/errors"
)

const (
	envoyYamlURL         = "https://raw.githubusercontent.com/solo-io/gloo/master/install/docker-compose/envoy-config.yaml"
	dockerComposeYamlURL = "https://raw.githubusercontent.com/solo-io/gloo/master/install/docker-compose/docker-compose.yaml"
)

func Install(folder string) error {
	err := createInstallFolder(folder)
	if err != nil {
		return err
	}
	err = download(dockerComposeYamlURL, filepath.Join(folder, "docker-compose.yaml"))
	if err != nil {
		return err
	}
	err = download(envoyYamlURL, filepath.Join(folder, "envoy-config.yaml"))
	if err != nil {
		return err
	}

	err = createStorageFolders(folder)
	if err != nil {
		return err
	}

	return updateGlooctlConfig(folder)
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

func createStorageFolders(folder string) error {
	for _, f := range []string{"_gloo_secrets", "_gloo_files"} {
		err := os.MkdirAll(filepath.Join(folder, f), 0755)
		if err != nil {
			return errors.Wrap(err, "unable to create storage directory "+f)
		}
	}

	// _gloo_config/*
	for _, f := range []string{"upstreams", "virtualservices"} {
		err := os.MkdirAll(filepath.Join(folder, "_gloo_config", f), 0755)
		if err != nil {
			return errors.Wrap(err, "unable to create storage directory"+f)
		}
	}

	return nil
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

func download(src, dst string) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
