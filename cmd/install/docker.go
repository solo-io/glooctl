package install

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/solo-io/glooctl/pkg/install/docker"
	"github.com/spf13/cobra"
)

const (
	successMessage = `Gloo setup successfully.
Please switch to directory '%s', and run "docker-compose up"
to start gloo.

`
)

func dockerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker [folder]",
		Short: "install gloo with Docker and file-based storage",
		Long: `
Installs gloo to run with Docker Compose in the given install folder.
If the folder doesn't exist glooctl will create it.

Once installed you can go to the install folder and run:
	docker-compose up 
to start gloo.

Glooctl will configure itself to use this instance of gloo.`,
		Args: cobra.ExactArgs(1),
		Run: func(c *cobra.Command, a []string) {
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Unable to get current directory", err)
				os.Exit(1)
			}
			installDir := filepath.Join(pwd, a[0])
			err = docker.Install(installDir)
			if err != nil {
				fmt.Printf("Unable to install gloo to %s: %q\n", installDir, err)
				os.Exit(1)
			}
			fmt.Printf(successMessage, installDir)
		},
	}
	return cmd
}
