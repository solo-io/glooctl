package virtualservice_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

func TestVirtualService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Virtual Service Suite")
}

var _ = BeforeSuite(helper.Build)
var _ = AfterSuite(helper.CleanUp)
