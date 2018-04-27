package upstream_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Creating upstream", func() {
	BeforeEach(setupStorage)
	AfterEach(tearDownStorage)

	It("should exit with exit code 1 when creating invalid upstream", func() {
		opts := withStorageOpts("upstream", "create", "-f", "testdata/invalid.yaml")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("missing secret reference"))
		Eventually(session).Should(gexec.Exit(1))
	})

	It("should create a valid upstream", func() {
		createAWSSecret()

		opts := withStorageOpts("upstream", "create", "-f", "testdata/aws.yaml")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		// check by doing a get
		opts = []string{"upstream", "get", "-o template", "--template={{range .}}{{.Name}} {{end}}"}
		opts = append(opts, storageOpts...)
		command = exec.Command(glooctlBinary, opts...)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(session.Wait().Out.Contents()).Should(ContainSubstring("testupstream"))
	})
})
