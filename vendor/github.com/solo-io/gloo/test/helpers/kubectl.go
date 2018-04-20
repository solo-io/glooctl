package helpers

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"os"

	"fmt"
	"strings"
	"time"

	"bytes"
	"text/template"

	"github.com/onsi/ginkgo"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/log"
)

const (
	// gloo labels
	testrunner            = "testrunner"
	helloservice          = "helloservice"
	helloservice2         = "helloservice-2"
	envoy                 = "ingress"
	gloo                  = "control-plane"
	kubeIngressController = "kube-ingress-controller"
	upstreamDiscovery     = "upstream-discovery"
	funcitonDiscovery     = "function-discovery"
	upstreamForEvents     = "upstream-for-events"
	grpcTestService       = "grpc-test-service"
	eventEmitter          = "event-emitter"
)

func SetupKubeForTest(namespace string) error {
	context := os.Getenv("KUBECTL_CONTEXT")
	if context == "" {
		current, err := KubectlOut("config", "current-context")
		if err != nil {
			return errors.Wrap(err, "getting currrent context")
		}
		context = strings.TrimSuffix(current, "\n")
	}
	// TODO(yuval-k): this changes the context for the user? can we do this less intrusive? maybe add it to
	// each kubectl command?
	if err := kubectl("config", "set-context", context, "--namespace="+namespace); err != nil {
		return errors.Wrap(err, "setting context")
	}
	return kubectl("create", "namespace", namespace)
}

func TeardownKube(namespace string) error {
	return kubectl("delete", "namespace", namespace)
}
func TeardownKubeE2E(namespace string) error {
	TeardownKube(namespace)
	kubectl("delete", "-f", filepath.Join(KubeE2eDirectory(), "kube_resources", "install.yml"))
	return kubectl("delete", "-f", filepath.Join(KubeE2eDirectory(), "kube_resources", "testing-resources.yml"))
}

func SetupKubeForE2eTest(namespace string, buildImages, push, debug bool) error {
	if err := SetupKubeForTest(namespace); err != nil {
		return err
	}
	if buildImages {
		if err := BuildPushContainers(push, debug); err != nil {
			return err
		}
	}
	kubeResourcesDir := filepath.Join(KubeE2eDirectory(), "kube_resources")

	envoyImageTag := os.Getenv("ENVOY_IMAGE_TAG")
	if envoyImageTag == "" {
		log.Warnf("no ENVOY_IMAGE_TAG specified, defaulting to latest")
		envoyImageTag = "latest"
	}

	pullPolicy := "IfNotPresent"

	if push {
		pullPolicy = "Always"
	}

	data := &struct {
		Namespace       string
		ImageTag        string
		ImagePullPolicy string
		Debug           string
	}{Namespace: namespace, ImageTag: ImageTag(), ImagePullPolicy: pullPolicy, Debug: ""}
	if debug {
		data.Debug = "-debug"
	}

	tmpl, err := template.New("Test_Resources").ParseFiles(filepath.Join(kubeResourcesDir, "helm-values.yaml.tmpl"))
	if err != nil {
		return errors.Wrap(err, "parsing template from helm-values.yaml.tmpl")
	}

	buf := &bytes.Buffer{}
	if err := tmpl.ExecuteTemplate(buf, "helm-values.yaml.tmpl", data); err != nil {
		return errors.Wrap(err, "executing template")
	}

	err = ioutil.WriteFile(filepath.Join(kubeResourcesDir, "helm-values.yaml"), buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, "writing generated test resources template")
	}

	installBytes, err := exec.Command("helm", "template", HelmDirectory(),
		"--namespace", namespace,
		"-n", "test",
		"--set", "ingress.imageTag="+envoyImageTag,
		"-f", filepath.Join(kubeResourcesDir, "helm-values.yaml")).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "running helm template: %v", string(installBytes))
	}

	err = ioutil.WriteFile(filepath.Join(kubeResourcesDir, "install.yml"), installBytes, 0644)
	if err != nil {
		return errors.Wrap(err, "writing generated install template")
	}

	tmpl, err = template.New("Test_Resources").ParseFiles(filepath.Join(kubeResourcesDir, "testing-resources.yaml.tmpl"))
	if err != nil {
		return errors.Wrap(err, "parsing template from testing-resources.yaml.tmpl")
	}

	buf = &bytes.Buffer{}
	if err := tmpl.ExecuteTemplate(buf, "testing-resources.yaml.tmpl", data); err != nil {
		return errors.Wrap(err, "executing template")
	}

	err = ioutil.WriteFile(filepath.Join(kubeResourcesDir, "testing-resources.yml"), buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, "writing generated test resources template")
	}

	// test stuff first
	if err := kubectl("apply", "-f", filepath.Join(kubeResourcesDir, "testing-resources.yml")); err != nil {
		return errors.Wrapf(err, "creating kube resource from testing-resources.yml")
	}
	if err := waitPodsRunning(
		testrunner,
		helloservice,
		helloservice2,
		upstreamForEvents,
		grpcTestService,
		eventEmitter,
	); err != nil {
		return errors.Wrap(err, "waiting for pods to start")
	}

	if err := kubectl("apply", "-f", filepath.Join(kubeResourcesDir, "install.yml")); err != nil {
		return errors.Wrapf(err, "creating kube resource from install.yml")
	}
	if err := waitPodsRunning(
		envoy,
		gloo,
		kubeIngressController,
		upstreamDiscovery,
		funcitonDiscovery,
	); err != nil {
		return errors.Wrap(err, "waiting for pods to start")
	}
	TestRunner("curl", "test-ingress:19000/logging?config=debug")
	return nil
}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	log.Debugf("k command: %v", cmd.Args)
	cmd.Env = os.Environ()
	// disable DEBUG=1 from getting through to kube
	for i, pair := range cmd.Env {
		if strings.HasPrefix(pair, "DEBUG") {
			cmd.Env = append(cmd.Env[:i], cmd.Env[i+1:]...)
			break
		}
	}
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	return cmd.Run()
}

func KubectlOut(args ...string) (string, error) {
	cmd := exec.Command("kubectl", args...)
	cmd.Env = os.Environ()
	// disable DEBUG=1 from getting through to kube
	for i, pair := range cmd.Env {
		if strings.HasPrefix(pair, "DEBUG") {
			cmd.Env = append(cmd.Env[:i], cmd.Env[i+1:]...)
			break
		}
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s (%v)", out, err)
	}
	return string(out), err
}

func KubectlOutAsync(args ...string) (*bytes.Buffer, chan struct{}, error) {
	cmd := exec.Command("kubectl", args...)
	cmd.Env = os.Environ()
	// disable DEBUG=1 from getting through to kube
	for i, pair := range cmd.Env {
		if strings.HasPrefix(pair, "DEBUG") {
			cmd.Env = append(cmd.Env[:i], cmd.Env[i+1:]...)
			break
		}
	}
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Start()
	if err != nil {
		err = fmt.Errorf("%s (%v)", buf.Bytes(), err)
	}
	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			cmd.Process.Kill()
		}
	}()
	return buf, done, err
}

// waitPodsRunning waits for all pods to be running
func waitPodsRunning(podNames ...string) error {
	finished := func(output string) bool {
		return strings.Contains(output, "Running")
	}
	for _, pod := range podNames {
		if err := waitPodStatus(pod, "Running", finished); err != nil {
			return err
		}
	}
	return nil
}

// waitPodsTerminated waits for all pods to be terminated
func waitPodsTerminated(podNames ...string) error {
	for _, pod := range podNames {
		finished := func(output string) bool {
			return !strings.Contains(output, pod)
		}
		if err := waitPodStatus(pod, "terminated", finished); err != nil {
			return err
		}
	}
	return nil
}

// TestRunner executes a command inside the TestRunner container
func TestRunner(args ...string) (string, error) {
	args = append([]string{"exec", "-i", testrunner, "--"}, args...)
	return KubectlOut(args...)
}

// TestRunnerAsync executes a command inside the TestRunner container
// returning a buffer that can be read from as it executes
func TestRunnerAsync(args ...string) (*bytes.Buffer, chan struct{}, error) {
	args = append([]string{"exec", "-i", testrunner, "--"}, args...)
	return KubectlOutAsync(args...)
}

func waitPodStatus(pod, status string, finished func(output string) bool) error {
	timeout := time.Second * 20
	interval := time.Millisecond * 1000
	tick := time.Tick(interval)

	log.Debugf("waiting %v for pod %v to be %v...", timeout, pod, status)
	for {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("timed out waiting for %v to be %v", pod, status)
		case <-tick:
			out, err := KubectlOut("get", "pod", "-l", "gloo="+pod)
			if err != nil {
				return fmt.Errorf("failed getting pod: %v", err)
			}
			if strings.Contains(out, "CrashLoopBackOff") {
				out = KubeLogs(pod)
				return errors.Errorf("%v in crash loop with logs %v", pod, out)
			}
			if strings.Contains(out, "ErrImagePull") || strings.Contains(out, "ImagePullBackOff") {
				out, _ = KubectlOut("describe", "pod", "-l", "gloo="+pod)
				return errors.Errorf("%v in ErrImagePull with description %v", pod, out)
			}
			if finished(out) {
				return nil
			}
		}
	}
}

func KubeLogs(pod string) string {
	out, err := KubectlOut("logs", "-l", "gloo="+pod)
	if err != nil {
		out = err.Error()
	}
	return out
}

func waitNamespaceStatus(namespace, status string, finished func(output string) bool) error {
	timeout := time.Second * 20
	interval := time.Millisecond * 1000
	tick := time.Tick(interval)

	log.Debugf("waiting %v for namespace %v to be %v...", timeout, namespace, status)
	for {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("timed out waiting for %v to be %v", namespace, status)
		case <-tick:
			out, err := KubectlOut("get", "namespace", namespace)
			if err != nil {
				return fmt.Errorf("failed getting pod: %v", err)
			}
			if finished(out) {
				return nil
			}
		}
	}
}
