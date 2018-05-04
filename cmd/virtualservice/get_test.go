package virtualservice_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Getting virtual service", func() {
	BeforeEach(setupVirtualServices)
	AfterEach(helper.TearDownStorage)

	It("should get list of virtual services for JSON output", func() {
		helper.RunWithArgs("virtualservice", "get", "-o", "json").
			ExpectExitCodeAndOutput(0, `"name":"axhixh.com"`)
	})

	It("should get list of virtual services for YAML output", func() {
		helper.RunWithArgs("virtualservice", "get", "-o", "yaml").
			ExpectExitCodeAndOutput(0, `name: axhixh.com`)
	})

	It("should get list of virtual services for template output", func() {
		helper.RunWithArgs("virtualservice", "get", "-o", "template",
			"--template", "{{range .}}{{.Name}} {{end}}").
			ExpectExitCodeAndOutput(0, "axhixh.com")
	})

	It("shouldget list of virtual services for default table output", func() {
		helper.RunWithArgs("virtualservice", "get").
			ExpectExitCodeAndOutput(0, "VIRTUAL SERVICE", "| axhixh.com ")
	})

	It("should get specific virtual service if a name is given", func() {
		helper.RunWithArgs("virtualservice", "get", "axhixh.com").
			ExpectExitCodeAndOutput(0, "| axhixh.com ")
	})

	It("should exit with satus code 1 if a name of invalid virtual service is given", func() {
		helper.RunWithArgs("virtualservice", "get", "non-exist").
			ExpectExitCodeAndOutput(1, "Unable to get virtual service")
	})
})
