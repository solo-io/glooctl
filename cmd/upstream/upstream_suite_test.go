package upstream_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/onsi/gomega/gexec"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/secretstorage"
	"github.com/solo-io/gloo/pkg/plugins/aws"
	"github.com/solo-io/gloo/pkg/storage/dependencies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	glooctlBinary string

	tmpDir    string
	configDir string
	secretDir string

	storageOpts []string
)

func TestUpstreams(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Upstream Suite")
}

var _ = BeforeSuite(func() {
	var err error
	glooctlBinary, err = gexec.Build("github.com/solo-io/glooctl")
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func setupStorage() {
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

func tearDownStorage() {
	err := os.RemoveAll(tmpDir)
	Expect(err).NotTo(HaveOccurred())
}

func withStorageOpts(opts ...string) []string {
	return append(opts, storageOpts...)
}
func bootstrapOpts() *bootstrap.Options {
	opts := &bootstrap.Options{}
	opts.ConfigStorageOptions.Type = "file"
	opts.SecretStorageOptions.Type = "file"
	opts.FileOptions.ConfigDir = configDir
	opts.FileOptions.SecretDir = secretDir

	return opts
}

func createAWSSecret() {
	opts := bootstrapOpts()
	ss, err := secretstorage.Bootstrap(*opts)
	Expect(err).NotTo(HaveOccurred())

	_, err = ss.Create(&dependencies.Secret{
		Ref:  "aws-secret.yaml",
		Data: map[string]string{aws.AwsAccessKey: "xxxxx", aws.AwsSecretKey: "yyyy"},
	})
	Expect(err).NotTo(HaveOccurred())
}

func setupUpstreams() {
	setupStorage()
	createAWSSecret()

	for _, f := range []string{"testdata/basic.yaml", "testdata/withfunctions.yaml"} {
		opts := withStorageOpts("upstream", "create", "-f", f)
		command := exec.Command(glooctlBinary, opts...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
	}

}
