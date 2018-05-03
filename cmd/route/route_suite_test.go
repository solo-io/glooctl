package route_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Route Suite")
}

var _ = BeforeSuite(helper.Build)
var _ = AfterSuite(helper.CleanUp)
