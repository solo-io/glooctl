package upstream_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Getting upstream", func() {
	BeforeEach(setupUpstreams)
	AfterEach(helper.TearDownStorage)

	It("should exit with exit code 1 when updating non existing upstream", func() {
		helper.RunWithArgs("upstream", "update", "-f", "testdata/update-non-exist.yaml").
			ExpectExitCodeAndOutput(1, `unable to find existing`)
	})
	It("should update valid upstream", func() {
		helper.RunWithArgs("upstream", "update", "-f", "testdata/update.yaml").ExpectExitCode(0)

		// verify update
		helper.RunWithArgs("upstream", "get", "with-function", "-o", "yaml").
			ExpectExitCodeAndOutput(0, `region: us-west-2`)
	})
})
