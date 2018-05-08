package openshift

import (
	"os"
	"os/exec"
)

// Install setups Gloo on OpenShift using oc and current context
func Install(dryRun bool) error {
	// running oc with latest install.yaml
	args := []string{"apply", "--filename",
		"https://raw.githubusercontent.com/solo-io/gloo/master/install/openshift/install.yaml"}
	if dryRun {
		args = append(args, "--dry-run=true")
	}
	cmd := exec.Command("oc", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
