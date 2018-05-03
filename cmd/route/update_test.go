package route_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Update routes", func() {
	BeforeEach(setupRoutes)
	AfterEach(helper.TearDownStorage)

	It("when specified matcher doesn't match existing routes", func() {
		helper.RunWithArgs("route", "update",
			"--old-path-exact", "/exact", "--old-upstream", "non-there", "--path-exact", "/exact2",
			"--upstream", "exact-upstream").
			ExpectExitCodeAndOutput(1, "could not find a route for the specified matcher and destination")
	})

	It("when specified matcher matches exactly one route", func() {
		helper.RunWithArgs("route", "update",
			"--old-path-exact", "/exact", "--old-upstream", "exact-upstream", "--path-exact", "/exact2",
			"--upstream", "exact-upstream").
			ExpectExitCodeAndOutput(0, "/exact2", "exact-upstream")
	})

})
