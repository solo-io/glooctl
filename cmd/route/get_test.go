package route_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Getting routes", func() {
	BeforeEach(setupRoutes)
	AfterEach(helper.TearDownStorage)

	It("without specifying output format", func() {
		helper.RunWithArgs("route", "get").
			ExpectExitCodeAndOutput(0,
				"POST", "func1", "GET", "test-upstream", "/exact")
	})

	It("when output format is YAML", func() {
		helper.RunWithArgs("route", "get", "--output", "yaml").
			ExpectExitCodeAndOutput(0, "path_prefix: /foo", "function_name: func1",
				"name: test-upstream", "path_exact: /exact")
	})

})
