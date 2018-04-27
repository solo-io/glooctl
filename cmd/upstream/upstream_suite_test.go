package upstream_test

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/solo-io/gloo/pkg/bootstrap/secretstorage"
	"github.com/solo-io/gloo/pkg/plugins/aws"
	"github.com/solo-io/gloo/pkg/storage/dependencies"
	"github.com/solo-io/glooctl/internal/test-helper"
)

func TestUpstreams(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Upstream Suite")
}

var _ = BeforeSuite(helper.Build)
var _ = AfterSuite(helper.CleanUp)

func createAWSSecret() {
	opts := helper.BootstrapOpts()
	ss, err := secretstorage.Bootstrap(*opts)
	Expect(err).NotTo(HaveOccurred())

	_, err = ss.Create(&dependencies.Secret{
		Ref:  "aws-secret.yaml",
		Data: map[string]string{aws.AwsAccessKey: "xxxxx", aws.AwsSecretKey: "yyyy"},
	})
	Expect(err).NotTo(HaveOccurred())
}

func setupUpstreams() {
	helper.SetupStorage()
	createAWSSecret()

	for _, f := range []string{"testdata/basic.yaml", "testdata/withfunctions.yaml"} {
		opts := helper.WithStorageOpts("upstream", "create", "-f", f)
		command := exec.Command(helper.Glooctl, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
	}

}
