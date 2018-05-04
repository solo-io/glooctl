package virtualservice_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Creating virtual service", func() {
	BeforeEach(helper.SetupStorage)
	AfterEach(helper.TearDownStorage)

	It("should allow creating a default virtual service when there isn't one", func() {
		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/mydefault.yaml").
			ExpectExitCodeAndOutput(0, "VIRTUAL SERVICE", "mydefault")
	})

	It("should not allow creating a new default virtual service if one exists", func() {
		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/mydefault.yaml").
			ExpectExitCode(0)

		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/vservice.yaml").
			ExpectExitCodeAndOutput(1, "domain")
	})

	It("should allow if it is a non default virtual service", func() {
		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/mydefault.yaml").
			ExpectExitCode(0)

		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/with-domains.yaml").
			ExpectExitCode(0)
		helper.RunWithArgs("virtualservice", "get", "axhixh.com", "-o", "yaml").
			ExpectExitCodeAndOutput(0, "domains:\n- axhixh.com")
	})
})
