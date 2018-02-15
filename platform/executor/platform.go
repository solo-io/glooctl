package executor

import (
	"fmt"
	"os"
	"path"
	"strings"

	storage "github.com/solo-io/glue-storage/pkg/storage"
	"github.com/solo-io/glue-storage/pkg/storage/crd"
	"github.com/solo-io/glue-storage/pkg/storage/file"
	"github.com/solo-io/gluectl/platform"
	"github.com/spf13/viper"
)

var (
	validStorage = map[string]string{
		"kubernetes": "file,crd",
	}
)

func GetExecutor(objname, namespace string) platform.Executor {
	// Read type from config and create executor for appropriate platform with some config args
	plat := viper.GetString("platform")

	if validStorage[plat] == "" {
		plat = "kubernetes"
	}

	stor := viper.GetString("storage")
	if stor == "" {
		stor = "crd"
	}

	if !strings.Contains(validStorage[plat], stor) {
		fmt.Printf("Invalid storage %s for platform %s\n", stor, plat)
		os.Exit(1)
	}
	var dataStore storage.Storage
	switch plat + "-" + stor {
	case "kubernetes-crd":
		kc := viper.GetString("kubeConfig")
		cfg, err := getClientConfig(kc)
		if err != nil {
			Fatal("Cannot create k8s client", err)
		}
		dataStore, err = crd.NewCrdStorage(cfg, namespace)
		if err != nil {
			Fatal("Cannot create crd storage", err)
		}

	case "kubernetes-file":
		root := viper.GetString("root")
		if root == "" {
			root = path.Join(os.Getenv("HOME"), ".glue")
		}
		var err error
		dataStore, err = file.NewFileStorage(root, namespace)
		if err != nil {
			Fatal("Cannot create file storage", err)
		}
	default:
		fmt.Printf("Cannot use platform %s with storage %s", plat, stor)
		os.Exit(1)
	}
	switch objname {
	case "upstream":
		return NewUpstreamExecutor(dataStore)
	case "vhost":
		return NewVhostExecutor(dataStore)
	}
	Fatal("Bad object name:", objname)
	return nil
}
