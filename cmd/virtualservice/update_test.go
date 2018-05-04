package virtualservice_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Updating virtual service", func() {
	BeforeEach(helper.SetupStorage)
	AfterEach(helper.TearDownStorage)

	It("should exist with code 1 when updating non existing virtual service", func() {
		helper.RunWithArgs("virtualservice", "update", "-f", "testdata/with-domains.yaml").
			ExpectExitCodeAndOutput(1, "unable to find existing virtual service axhixh.com")
	})

	It("should update valid virtual service", func() {
		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/with-domains.yaml").
			ExpectExitCodeAndOutput(0, "axhixh.com, www.axhixh.com")
		helper.RunWithArgs("virtualservice", "update", "-f", "testdata/with-domains-update.yaml").
			ExpectExitCodeAndOutput(0, "axhixh.net, www.axhixh.net")
	})
})
