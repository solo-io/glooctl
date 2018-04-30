package secret_test

import (
	. "github.com/onsi/ginkgo"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Creating secret", func() {
	BeforeEach(helper.SetupStorage)
	AfterEach(helper.TearDownStorage)

	It("should exit with exit code 1 when creating a Google secret with invalid file", func() {
		helper.RunWithArgs("secret", "create", "google", "--name=google-secret", "--filename=doesntexist").
			ExpectExitCodeAndOutput(1, "unable to read service account key file")
	})
})
