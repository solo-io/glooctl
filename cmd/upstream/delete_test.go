package upstream_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Deleting upstream", func() {
	BeforeEach(setupStorage)
	AfterEach(tearDownStorage)

	It("should exit with exit code 1 when deleting non existing upstream", func() {
		opts := withStorageOpts("upstream", "delete", "nonexist")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("Unable to delete upstream nonexist"))
		Eventually(session).Should(gexec.Exit(1))
	})

	It("should exist with exit code 1 when calling delete without upstream name", func() {
		opts := withStorageOpts("upstream", "delete")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(1))
	})

	It("should delete the upstream with given name", func() {
		// create
		opts := withStorageOpts("upstream", "create", "-f", "testdata/basic.yaml")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		// delete
		opts = withStorageOpts("upstream", "delete", "testupstream")
		command = exec.Command(glooctlBinary, opts...)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("Upstream testupstream deleted"))
		Eventually(session).Should(gexec.Exit(0))
	})
})
