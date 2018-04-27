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

	It("should get list of upstreams for JSON output", func() {
		opts := withStorageOpts("upstream", "get", "-o", "json")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say(`"name":"testupstream"`))
		Eventually(session.Out).Should(gbytes.Say(`"name":"with-function"`))
		Eventually(session).Should(gexec.Exit(0))
	})

	It("should get list of upstreams for YAML output", func() {
		opts := withStorageOpts("upstream", "get", "-o", "yaml")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say(`name: testupstream`))
		Eventually(session.Out).Should(gbytes.Say(`name: with-function`))
		Eventually(session).Should(gexec.Exit(0))
	})

	It("should get list of upstreams for template output", func() {
		opts := withStorageOpts("upstream", "get", "-o", "template", "--template", "{{range .}}{{.Name}} {{end}}")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say(`testupstream with-function`))
		Eventually(session).Should(gexec.Exit(0))
	})

	It("should get list of upstreams for default table output", func() {
		opts := withStorageOpts("upstream", "get")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say(`| testupstream`))
		Eventually(session.Out).Should(gbytes.Say(`| with-function`))
		Eventually(session.Out).Should(gbytes.Say(`| gloo-hello`))
		Eventually(session).Should(gexec.Exit(0))
	})

	It("should get specific upstream if a name is given", func() {
		opts := withStorageOpts("upstream", "get", "testupstream")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say(`| testupstream`))
		Eventually(session).Should(gexec.Exit(0))
	})

	It("should exit with status code 1 if a name of invalid upstream is given", func() {
		opts := withStorageOpts("upstream", "get", "non-exist")
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say(`Unable to get upstream`))
		Eventually(session).Should(gexec.Exit(1))
	})
})
