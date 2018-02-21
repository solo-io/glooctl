package util

import (
	"fmt"
	"time"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	"github.com/solo-io/gloo-storage/file"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

func GetStorageClient(c *cobra.Command) (storage.Interface, error) {
	flags := c.InheritedFlags()
	kubeConfig, _ := flags.GetString("kubeconfig")
	namespace, _ := flags.GetString("namespace")
	period, _ := flags.GetInt("sync-period")
	syncPeriod := time.Duration(period) * time.Second
	if kubeConfig != "" {
		kubeClient, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}
		return crd.NewStorage(kubeClient, namespace, syncPeriod)
	}
	glooFolder, _ := flags.GetString("gloo-folder")
	if glooFolder != "" {
		return file.NewStorage(glooFolder, syncPeriod)
	}
	return nil, fmt.Errorf("unable to create storage client")
}
