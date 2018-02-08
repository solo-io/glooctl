package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func GetGlobalFlags(cmd *cobra.Command) (string, string, int) {
	file, err := cmd.InheritedFlags().GetString("file")
	if err != nil {
		log.Fatal("Invalid value of the 'file' flag", err)
	}
	ns, err := cmd.InheritedFlags().GetString("namespace")
	if err != nil {
		log.Fatal("Invalid value of the 'namespace' flag", err)
	}
	wait, err := cmd.InheritedFlags().GetInt("wait")
	if err != nil {
		log.Fatal("Invalid value of the 'wait' flag", err)
	}
	return file, ns, wait
}
