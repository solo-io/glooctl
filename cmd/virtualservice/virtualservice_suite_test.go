package virtualservice_test

import (
	"bufio"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

func TestVirtualService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Virtual Service Suite")
}

var _ = BeforeSuite(helper.Build)
var _ = AfterSuite(helper.CleanUp)

func setupVirtualServices() {
	helper.SetupStorage()
	helper.RunWithArgs("virtualservice", "create", "-f", "testdata/mydefault.yaml").
		ExpectExitCode(0)
	helper.RunWithArgs("virtualservice", "create", "-f", "testdata/with-domains.yaml").
		ExpectExitCode(0)

}

// for interactive tests
func expect(expected string, buf *bufio.Reader) {
	helper.ExpectOutput(buf, expected)
}
