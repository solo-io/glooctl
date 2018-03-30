package secrets

import (
	"path/filepath"

	"github.com/solo-io/gloo-secret"
	"github.com/solo-io/gloo-secret/crd"
	"github.com/solo-io/gloo-secret/file"
	"github.com/solo-io/glooctl/pkg/util"
	"k8s.io/client-go/tools/clientcmd"
)

func GetSecretClient(opts *util.StorageOptions) (secret.SecretInterface, error) {
	secretDir := opts.SecretDir
	if secretDir != "" {
		return file.NewClient(secretDir)
	}

	kubeConfig := opts.KubeConfig
	if kubeConfig == "" && util.HomeDir() != "" {
		kubeConfig = filepath.Join(util.HomeDir(), ".kube", "config")
	}
	kubeClient, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	return crd.NewClient(kubeClient, opts.Namespace)
}
