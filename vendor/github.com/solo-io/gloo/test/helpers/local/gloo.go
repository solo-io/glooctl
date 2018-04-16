package localhelpers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/gloo/pkg/storage/file"

	"github.com/solo-io/gloo/test/helpers"

	"github.com/onsi/ginkgo"

	"github.com/ghodss/yaml"
)

type GlooFactory struct {
	srcpath  string
	gloopath string
	wasbuilt bool
}

func NewGlooFactory() (*GlooFactory, error) {
	gloopath := os.Getenv("GLOO_BINARY")

	if gloopath != "" {
		return &GlooFactory{
			gloopath: gloopath,
		}, nil
	}
	srcpath := filepath.Join(helpers.GlooSoloDirectory(), "cmd", "control-plane")
	gf := &GlooFactory{
		srcpath: srcpath,
	}
	err := gf.build()
	if err != nil {
		return nil, err
	}
	gloopath = filepath.Join(srcpath, "control-plane")
	gf.gloopath = gloopath
	return gf, nil
}

func (gf *GlooFactory) build() error {
	if gf.srcpath == "" {
		if gf.gloopath == "" {
			return errors.New("can't build gloo and none provided")
		}
		return nil
	}
	gf.wasbuilt = true

	cmd := exec.Command("go", "build", "-v", "-i", "-gcflags", "-N -l", "-o", "control-plane", "main.go")

	cmd.Dir = gf.srcpath
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (gf *GlooFactory) NewGlooInstance() (*GlooInstance, error) {

	tmpdir, err := ioutil.TempDir(os.Getenv("HELPER_TMP"), "gloo")
	if err != nil {
		return nil, err
	}

	gi := &GlooInstance{
		gloopath: gf.gloopath,
		tmpdir:   tmpdir,
	}

	if err := gi.initStorage(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return gi, nil
}

func (gf *GlooFactory) Clean() error {
	return nil
}

type GlooInstance struct {
	gloopath string

	tmpdir string
	store  storage.Interface
	cmd    *exec.Cmd
}

func (gi *GlooInstance) ConfigDir() string {
	return filepath.Join(gi.tmpdir)
}

func (gi *GlooInstance) FilesDir() string {
	return filepath.Join(gi.tmpdir, "_gloo_files")
}

func (gi *GlooInstance) SecretsDir() string {
	return filepath.Join(gi.tmpdir, "_gloo_secrets")
}

func (gi *GlooInstance) EnvoyPort() uint32 {
	return 8080
}

func (gi *GlooInstance) AddUpstream(u *v1.Upstream) error {
	_, err := gi.store.V1().Upstreams().Create(u)
	return err
}

func (gi *GlooInstance) GetUpstream(s string) (*v1.Upstream, error) {
	return gi.store.V1().Upstreams().Get(s)
}

func (gi *GlooInstance) AddVhost(u *v1.VirtualHost) error {
	_, err := gi.store.V1().VirtualHosts().Create(u)
	return err
}

func (gi *GlooInstance) AddSecret(name string, secret map[string]string) error {
	secretdir := filepath.Join(gi.tmpdir, "_gloo_secrets")
	os.Mkdir(secretdir, 0755)
	secretfile := filepath.Join(secretdir, name)

	data, err := yaml.Marshal(&secret)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(secretfile, data, 0400)
}

func (gi *GlooInstance) initStorage() error {
	dir := gi.tmpdir
	client, err := file.NewStorage(filepath.Join(dir, "_gloo_config"), time.Hour)
	if err != nil {
		return errors.New("failed to start file config watcher for directory " + dir)
	}
	err = client.V1().Register()
	if err != nil {
		return errors.New("failed to register file config watcher for directory " + dir)
	}
	// enable file storage
	if err := os.MkdirAll(filepath.Join(gi.tmpdir, "_gloo_files"), 0755); err != nil {
		return err
	}
	// enable secret storage
	if err := os.MkdirAll(filepath.Join(gi.tmpdir, "_gloo_secrets"), 0755); err != nil {
		return err
	}
	gi.store = client
	return nil

}
func (gi *GlooInstance) Run() error {
	return gi.RunWithPort(8081)
}

func (gi *GlooInstance) RunWithPort(xdsport uint32) error {

	var cmd *exec.Cmd
	glooargs := []string{
		"--storage.type=file",
		"--storage.refreshrate=1s",
		"--secrets.type=file",
		"--secrets.refreshrate=1s",
		fmt.Sprintf("--xds.port=%d", xdsport),
	}
	if os.Getenv("DEBUG_GLOO") != "" {
		dlvargs := append([]string{"--headless", "--listen=:2345", "--log", "exec", gi.gloopath, "--"}, glooargs...)
		cmd = exec.Command("dlv", dlvargs...)
	} else {
		cmd = exec.Command(gi.gloopath, glooargs...)
	}

	cmd.Dir = gi.tmpdir
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	err := cmd.Start()
	if err != nil {
		return err
	}
	gi.cmd = cmd
	return nil
}

func (gi *GlooInstance) Clean() error {
	if gi == nil {
		return nil
	}
	if gi.cmd != nil {
		gi.cmd.Process.Kill()
		gi.cmd.Wait()
	}
	if gi.tmpdir != "" {
		defer os.RemoveAll(gi.tmpdir)

	}

	return nil
}
