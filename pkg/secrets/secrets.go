package secrets

// FIXME - extract this out to its own repository for secret client repository

import (
	"path/filepath"

	"github.com/solo-io/gloo-secret"
	"github.com/solo-io/gloo-secret/crd"
	"github.com/solo-io/gloo-secret/file"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// FIXME(ashish) pass necessary parameters
func GetSecretClient(c *cobra.Command) (secret.SecretInterface, error) {
	flags := c.InheritedFlags()
	resourceFolder, _ := flags.GetString("resource-folder")
	if resourceFolder != "" {
		return file.NewClient(resourceFolder)
	}

	kubeConfig, _ := flags.GetString("kubeconfig")
	if kubeConfig == "" && util.HomeDir() != "" {
		kubeConfig = filepath.Join(util.HomeDir(), ".kube", "config")
	}
	kubeClient, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	namespace, _ := flags.GetString("namespace")
	return crd.NewClient(kubeClient, namespace)
}
