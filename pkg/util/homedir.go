package util

import (
	"os"
	"path/filepath"
)

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func ConfigDir() (string, error) {
	d := filepath.Join(HomeDir(), ".glooctl")
	_, err := os.Stat(d)
	if err == nil {
		return d, nil
	}
	if os.IsNotExist(err) {
		if err := os.Mkdir(d, 0755); err != nil {
			return "", err
		}
		return d, nil
	}

	return d, err
}
