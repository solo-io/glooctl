package upstream_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Getting upstream", func() {
	BeforeEach(setupUpstreams)
	AfterEach(tearDownStorage)

	It("should exit with exit code 1 when updating non existing upstream", func() {
		opts := withStorageOpts("upstream", "update", "-f", "testdata/update-non-exist.yaml")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say(`unable to find existing`))
		Eventually(session).Should(gexec.Exit(1))
	})
	It("should update valid upstream", func() {
		opts := withStorageOpts("upstream", "update", "-f", "testdata/update.yaml")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		// verify update
		opts = withStorageOpts("upstream", "get", "with-function", "-o", "yaml")
		command = exec.Command(glooctlBinary, opts...)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say(`region: us-west-2`))
		Eventually(session).Should(gexec.Exit(0))
	})
})
