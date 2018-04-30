package upstream_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Creating upstream", func() {
	BeforeEach(helper.SetupStorage)
	AfterEach(helper.TearDownStorage)

	It("should exit with exit code 1 when creating invalid upstream", func() {
		helper.RunWithArgs("upstream", "create", "-f", "testdata/invalid.yaml").
			ExpectExitCodeAndOutput(1, "missing secret reference")
	})

	It("should create a valid upstream", func() {
		createAWSSecret()
		helper.RunWithArgs("upstream", "create", "-f", "testdata/aws.yaml").ExpectExitCode(0)

		// check by doing a get
		helper.RunWithArgs("upstream", "get", "-o template", "--template={{range .}}{{.Name}} {{end}}").
			ExpectExitCodeAndOutput(0, "testupstream")
	})
})
