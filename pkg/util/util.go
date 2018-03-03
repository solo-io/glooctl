package util

import (
	"os"
	"path/filepath"
	"time"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	"github.com/solo-io/gloo-storage/file"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

func GetStorageClient(c *cobra.Command) (storage.Interface, error) {
	flags := c.InheritedFlags()
	period, _ := flags.GetInt("sync-period")
	syncPeriod := time.Duration(period) * time.Second

	resourceFolder, _ := flags.GetString("resource-folder")
	if resourceFolder != "" {
		return file.NewStorage(filepath.Join(resourceFolder, "config"), syncPeriod)
	}

	kubeConfig, _ := flags.GetString("kubeconfig")
	namespace, _ := flags.GetString("namespace")
	if kubeConfig == "" && HomeDir() != "" {
		kubeConfig = filepath.Join(HomeDir(), ".kube", "config")
	}
	kubeClient, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	return crd.NewStorage(kubeClient, namespace, syncPeriod)
}

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
