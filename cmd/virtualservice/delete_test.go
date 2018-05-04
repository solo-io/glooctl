package virtualservice_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Deleting virtual service", func() {
	BeforeEach(helper.SetupStorage)
	AfterEach(helper.TearDownStorage)

	It("should exit with code 1 when deleting virtual service that doesn't exist", func() {
		helper.RunWithArgs("virtualservice", "delete", "my-virtual-service").
			ExpectExitCodeAndOutput(1, "Unable to delete virtual service my-virtual-service")
	})

	It("should allow deleting a virtual service that exists", func() {
		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/with-domains.yaml").
			ExpectExitCode(0)
		helper.RunWithArgs("virtualservice", "delete", "axhixh.com").
			ExpectExitCodeAndOutput(0, "Virtual service axhixh.com deleted")
	})
})
