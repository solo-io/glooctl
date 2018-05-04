package kube

import (
	"os"
	"os/exec"
)

// Install setups Gloo on Kubernetes using kubectl and current context
func Install(dryRun bool) error {
	// using kubectl with latest install.yaml
	args := []string{"apply", "--filename",
		"https://raw.githubusercontent.com/solo-io/gloo/master/install/kube/install.yaml"}
	if dryRun {
		args = append(args, "--dry-run=true")
	}
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
