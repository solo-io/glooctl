package route_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Deleting route", func() {
	BeforeEach(setupRoutes)
	AfterEach(helper.TearDownStorage)

	It("when specified matcher doesn't match existing routes", func() {
		helper.RunWithArgs("route", "delete", "--path-exact", "/not-there",
			"--upstream", "test-upstream").
			ExpectExitCodeAndOutput(1, "did not match any route")
	})

	It("when specified matcher matches exactly one route", func() {
		helper.RunWithArgs("route", "delete", "--path-exact", "/exact",
			"--upstream", "exact-upstream").
			ExpectExitCode(0)
	})

})
