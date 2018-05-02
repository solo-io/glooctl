package kube_e2e

import (
	. "github.com/onsi/ginkgo"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/test/helpers"

	"os"
	"testing"
)

func TestKubernetes(t *testing.T) {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}

	helpers.RegisterPreFailHandler(func() {
		var logs string
		for _, component := range []string{"control-plane", "ingress"} {
			l, err := helpers.KubectlOut("logs", "-l", "gloo="+component)
			logs += string(l) + "\n"
			if err != nil {
				logs += "error getting logs for " + component + ": " + err.Error()
			}
		}

		log.Printf("\n****************************************" +
			"\nLOGS FROM THE KUBE BOYS: \n\n" + logs + "\n************************************")
	})

	helpers.RegisterCommonFailHandlers()
	log.DefaultOut = GinkgoWriter
	RunSpecs(t, "Kubernetes Suite")
}
