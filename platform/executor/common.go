package executor

import (
	"fmt"
	"os"

	"github.com/solo-io/gluectl/platform"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getClientConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
	return rest.InClusterConfig()
}

func Fatal(x ...interface{}) {
	fmt.Println("\nERROR: ", x)
	os.Exit(1)
}

func getUParams(params interface{}) *platform.UpstreamParams {
	return params.(*platform.UpstreamParams)
}

func getVParams(params interface{}) *platform.VhostParams {
	return params.(*platform.VhostParams)
}
