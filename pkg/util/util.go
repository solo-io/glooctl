package util

import (
	"os"
	"path/filepath"
	"time"

	"log"

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

	resourceFolder, _ := flags.GetString("gloo-config-dir")
	if resourceFolder != "" {
		log.Printf("Using file-based storage for gloo. Gloo must be configured to use file storage with config dir %v", resourceFolder)
		return file.NewStorage(resourceFolder, syncPeriod)
	}
	log.Printf("Using kubernetes crd-based storage for gloo. Gloo must be configured to use kubernetes storage")

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
