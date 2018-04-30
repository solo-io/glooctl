package upstream_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Deleting upstream", func() {
	BeforeEach(helper.SetupStorage)
	AfterEach(helper.TearDownStorage)

	It("should exit with exit code 1 when deleting non existing upstream", func() {
		helper.RunWithArgs("upstream", "delete", "nonexist").
			ExpectExitCodeAndOutput(1, "Unable to delete upstream nonexist")
	})

	It("should exist with exit code 1 when calling delete without upstream name", func() {
		helper.RunWithArgs("upstream", "delete").ExpectExitCode(1)
	})

	It("should delete the upstream with given name", func() {
		// create
		helper.RunWithArgs("upstream", "create", "-f", "testdata/basic.yaml").ExpectExitCode(0)

		// delete
		helper.RunWithArgs("upstream", "delete", "testupstream").
			ExpectExitCodeAndOutput(0, "Upstream testupstream deleted")
	})
})
