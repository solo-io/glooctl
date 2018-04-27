package upstream_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Deleting upstream", func() {
	BeforeEach(helper.SetupStorage)
	AfterEach(helper.TearDownStorage)

	It("should exit with exit code 1 when deleting non existing upstream", func() {
		opts := helper.WithStorageOpts("upstream", "delete", "nonexist")
		command := exec.Command(helper.Glooctl, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("Unable to delete upstream nonexist"))
		Eventually(session).Should(gexec.Exit(1))
	})

	It("should exist with exit code 1 when calling delete without upstream name", func() {
		opts := helper.WithStorageOpts("upstream", "delete")
		command := exec.Command(helper.Glooctl, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(1))
	})

	It("should delete the upstream with given name", func() {
		// create
		opts := helper.WithStorageOpts("upstream", "create", "-f", "testdata/basic.yaml")
		command := exec.Command(helper.Glooctl, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		// delete
		opts = helper.WithStorageOpts("upstream", "delete", "testupstream")
		command = exec.Command(helper.Glooctl, opts...)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("Upstream testupstream deleted"))
		Eventually(session).Should(gexec.Exit(0))
	})
})
