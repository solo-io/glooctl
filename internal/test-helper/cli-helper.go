package helper

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kr/pty"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/solo-io/gloo/pkg/bootstrap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	RED   = "\033[31m"
	RESET = "\033[0m"
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

	err = os.MkdirAll(filepath.Join(configDir, "virtualservices"), 0700)
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
	command.Env = append(os.Environ(), "CHECKPOINT_DISABLE=1")
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(code))
}

// ExpectExitCodeAndOutput runs glooctl and expects the given exit code and
// output message
func (a *Args) ExpectExitCodeAndOutput(code int, messages ...string) {
	command := exec.Command(Glooctl, a.Opts...)
	command.Env = append(os.Environ(), "CHECKPOINT_DISABLE=1")
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	for _, m := range messages {
		Eventually(session.Out).Should(gbytes.Say(m))
	}
	Eventually(session).Should(gexec.Exit(code))
}

func (a *Args) Interact(code int, interaction func(*bufio.Reader, *os.File)) {
	fh, tty, err := pty.Open()
	Ω(err).ShouldNot(HaveOccurred())
	defer tty.Close()
	defer fh.Close()

	command := exec.Command(Glooctl, a.Opts...)
	command.Env = append(os.Environ(), "CHECKPOINT_DISABLE=1")
	command.Stdin = tty
	session, err := gexec.Start(command, tty, tty)

	buf := bufio.NewReaderSize(fh, 1024)

	interaction(buf, fh)

	Eventually(session).Should(gexec.Exit(code))
}

// ExpectOutput compares the output of the interaction with expected
// Taken from generated code from autoplay
func ExpectOutput(buf *bufio.Reader, expected string) {
	sofar := []rune{}
	for _, r := range expected {
		got, _, _ := buf.ReadRune()
		sofar = append(sofar, got)
		if got != r {
			fmt.Fprintln(os.Stderr, RESET)

			// we want to quote the string but we also want to make the unexpected character RED
			// so we use the strconv.Quote function but trim off the quoted characters so we can
			// merge multiple quoted strings into one.
			expStart := strings.TrimSuffix(strconv.Quote(expected[:len(sofar)-1]), "\"")
			expMiss := strings.TrimSuffix(strings.TrimPrefix(strconv.Quote(string(expected[len(sofar)-1])), "\""), "\"")
			expEnd := strings.TrimPrefix(strconv.Quote(expected[len(sofar):]), "\"")

			fmt.Fprintf(os.Stderr, "Expected: %s%s%s%s%s\n", expStart, RED, expMiss, RESET, expEnd)

			// read the rest of the buffer
			p := make([]byte, buf.Buffered())
			buf.Read(p)

			gotStart := strings.TrimSuffix(strconv.Quote(string(sofar[:len(sofar)-1])), "\"")
			gotMiss := strings.TrimSuffix(strings.TrimPrefix(strconv.Quote(string(sofar[len(sofar)-1])), "\""), "\"")
			gotEnd := strings.TrimPrefix(strconv.Quote(string(p)), "\"")

			fmt.Fprintf(os.Stderr, "Got:      %s%s%s%s%s\n", gotStart, RED, gotMiss, RESET, gotEnd)
			Fail(fmt.Sprintf("Unexpected Rune %q, Expected %q\n", got, r))
		} else {
			fmt.Printf("%c", r)
		}
	}
}
