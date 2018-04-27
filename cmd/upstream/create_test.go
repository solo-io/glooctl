package upstream_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Creating upstream", func() {
	BeforeEach(helper.SetupStorage)
	AfterEach(helper.TearDownStorage)

	It("should exit with exit code 1 when creating invalid upstream", func() {
		opts := helper.WithStorageOpts("upstream", "create", "-f", "testdata/invalid.yaml")
		command := exec.Command(helper.Glooctl, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("missing secret reference"))
		Eventually(session).Should(gexec.Exit(1))
	})

	It("should create a valid upstream", func() {
		createAWSSecret()

		opts := helper.WithStorageOpts("upstream", "create", "-f", "testdata/aws.yaml")
		command := exec.Command(helper.Glooctl, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		// check by doing a get
		opts = helper.WithStorageOpts("upstream", "get", "-o template", "--template={{range .}}{{.Name}} {{end}}")
		command = exec.Command(helper.Glooctl, opts...)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(session.Wait().Out.Contents()).Should(ContainSubstring("testupstream"))
	})
})
