package secret_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/glooctl/internal/test-helper"
)

func TestSecrets(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Secret Suite")
}

var _ = BeforeSuite(helper.Build)
var _ = AfterSuite(helper.CleanUp)
