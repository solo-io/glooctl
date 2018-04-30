package upstream_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Getting upstream", func() {
	BeforeEach(setupUpstreams)
	AfterEach(helper.TearDownStorage)

	It("should get list of upstreams for JSON output", func() {
		helper.RunWithArgs("upstream", "get", "-o", "json").
			ExpectExitCodeAndOutput(0, `"name":"testupstream"`, `"name":"with-function"`)
	})

	It("should get list of upstreams for YAML output", func() {
		helper.RunWithArgs("upstream", "get", "-o", "yaml").
			ExpectExitCodeAndOutput(0, `name: testupstream`, `name: with-function`)
	})

	It("should get list of upstreams for template output", func() {
		helper.RunWithArgs("upstream", "get", "-o", "template", "--template", "{{range .}}{{.Name}} {{end}}").
			ExpectExitCodeAndOutput(0, `testupstream`, `with-function`)
	})

	It("should get list of upstreams for default table output", func() {
		helper.RunWithArgs("upstream", "get").
			ExpectExitCodeAndOutput(0, `| testupstream`, `| with-function`, `| gloo-hello`)
	})

	It("should get specific upstream if a name is given", func() {
		helper.RunWithArgs("upstream", "get", "testupstream").ExpectExitCodeAndOutput(0, `| testupstream`)
	})

	It("should exit with status code 1 if a name of invalid upstream is given", func() {
		helper.RunWithArgs("upstream", "get", "non-exist").ExpectExitCodeAndOutput(1, `Unable to get upstream`)
	})
})
