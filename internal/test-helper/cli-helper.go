package helper

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/solo-io/gloo/pkg/bootstrap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	// Glooctl points to the newly created binary
	Glooctl string

	tmpDir    string
	configDir string
	secretDir string

	storageOpts []string
)

// Build - builds the glooctl binary for testing
func Build() {
	var err error
	Glooctl, err = gexec.Build("github.com/solo-io/glooctl")
	Ω(err).ShouldNot(HaveOccurred())
}

// CleanUp - cleans any binaries created for test
func CleanUp() {
	gexec.CleanupBuildArtifacts()
}

// SetupStorage sets up file based storage for testing glooctl
func SetupStorage() {
	By("Creating temporary directory for file storage")

	var err error
	tmpDir, err = ioutil.TempDir("", "glooctl-test")
	Expect(err).NotTo(HaveOccurred())

	configDir = filepath.Join(tmpDir, "config")
	secretDir = filepath.Join(tmpDir, "secret")

	err = os.MkdirAll(filepath.Join(configDir, "upstreams"), 0700)
	Expect(err).NotTo(HaveOccurred())

	err = os.MkdirAll(secretDir, 0700)
	Expect(err).NotTo(HaveOccurred())

	storageOpts = []string{"--secrets.type=file",
		"--storage.type=file",
		"--file.config.dir=" + configDir,
		"--file.secret.dir=" + secretDir,
	}
}

// TearDownStorage - cleans up file based storage used for testing
func TearDownStorage() {
	err := os.RemoveAll(tmpDir)
	Expect(err).NotTo(HaveOccurred())
}

// BootstrapOpts returns the options used to represent the storage used
func BootstrapOpts() *bootstrap.Options {
	opts := &bootstrap.Options{}
	opts.ConfigStorageOptions.Type = "file"
	opts.SecretStorageOptions.Type = "file"
	opts.FileOptions.ConfigDir = configDir
	opts.FileOptions.SecretDir = secretDir

	return opts
}

type Args struct {
	Opts []string
}

// RunWithArgs setups glooctl to run with given CLI parameters
func RunWithArgs(opts ...string) *Args {
	return &Args{Opts: append(opts, storageOpts...)}
}

// ExpectExitCode runs glooctl and expects the given exit code
func (a *Args) ExpectExitCode(code int) {
	command := exec.Command(Glooctl, a.Opts...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(code))
}

// ExpectExitCodeAndOutput runs glooctl and expects the given exit code and
// output message
func (a *Args) ExpectExitCodeAndOutput(code int, messages ...string) {
	command := exec.Command(Glooctl, a.Opts...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	for _, m := range messages {
		Eventually(session.Out).Should(gbytes.Say(m))
	}
	Eventually(session).Should(gexec.Exit(code))
}
