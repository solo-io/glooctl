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

func setupRoutes() {
	helper.SetupStorage()

	helper.RunWithArgs("route", "create", "--path-prefix", "/foo", "--http-method", "POST",
		"--upstream", "upstream1", "--function", "func1").
		ExpectExitCode(0)
	helper.RunWithArgs("route", "create", "--path-prefix", "/foo", "--http-method", "GET",
		"--upstream", "test-upstream").
		ExpectExitCode(0)
	helper.RunWithArgs("route", "create", "--path-exact", "/exact", "--upstream", "exact-upstream").
		ExpectExitCode(0)
}
