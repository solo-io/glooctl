package secrets

// FIXME - extract this out to its own repository for secret client repository

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// FIXME(ashish) pass necessary parameters
func GetSecretClient(c *cobra.Command) (v1.SecretInterface, error) {
	flags := c.InheritedFlags()
	resourceFolder, _ := flags.GetString("resource-folder")
	if resourceFolder != "" {
		return nil, fmt.Errorf("File based secret client not implemented yet")
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
	cs, err := kubernetes.NewForConfig(kubeClient)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get kubernetes client")
	}
	return cs.CoreV1().Secrets(namespace), nil
}
