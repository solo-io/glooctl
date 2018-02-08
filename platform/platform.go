package platform

import (
	"github.com/solo-io/gluectl/platform/k8s"
	"github.com/spf13/viper"
)

func GetExecutor() Executor {
	// Read type from config and create executor for appropriate platform with some config args
	switch t := viper.GetString("platform"); t {
	default:
		kc := viper.GetString("kubeConfig")
		return k8s.NewExecutor(kc)
	}
}
