package cmd

import (
	"github.com/solo-io/gluectl/platform"
	"github.com/spf13/cobra"
)

var (
	gparams = &platform.GlobalParams{}
)

func CreateGlobalFlags(cmd *cobra.Command, isWait bool) {
	cmd.PersistentFlags().StringVarP(&gparams.FileName, "file", "f", "", "file with resource definition")
	cmd.PersistentFlags().StringVarP(&gparams.Namespace, "namespace", "n", "", "resource namespace")

	if isWait {
		cmd.PersistentFlags().IntVarP(&gparams.WaitSec, "wait", "w", 0, "seconds to wait, 0 - return immediately")
	}
}

func GetGlobalFlags() *platform.GlobalParams {
	return gparams
}
