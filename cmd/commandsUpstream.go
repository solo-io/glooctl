package cmd

import (
	"github.com/solo-io/gluectl/platform/common"
	"github.com/spf13/cobra"
)

var upstreamCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Create upstream",
	Long:  `Create upstream configuration object using config file and/or command line arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		LoadUpstreamParamsFromFile()
		ex := common.GetExecutor()
		gp := GetGlobalFlags()
		up := GetUpstreamParams()
		ex.RunCreateUpstream(gp, up)
	},
}

var upstreamDelCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Delete upstream",
	Long:  `Delete Upstream by name`,
	Run: func(cmd *cobra.Command, args []string) {
		LoadUpstreamParamsFromFile()
		ex := common.GetExecutor()
		gp := GetGlobalFlags()
		up := GetUpstreamParams()
		ex.RunDeleteUpstream(gp, up)
	},
}

var upstreamDescribeCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Describe upstream(s)",
	Long:  `Describe upstream (by name) or all upstreams in the namespace`,

	Run: func(cmd *cobra.Command, args []string) {
		LoadUpstreamParamsFromFile()
		ex := common.GetExecutor()
		gp := GetGlobalFlags()
		up := GetUpstreamParams()
		ex.RunDescribeUpstream(gp, up)
	},
}

var upstreamGetCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Get upstream(s)",
	Long:  `Get upstream (by name) or all upstreams in the namespace`,

	Run: func(cmd *cobra.Command, args []string) {
		LoadUpstreamParamsFromFile()
		ex := common.GetExecutor()
		gp := GetGlobalFlags()
		up := GetUpstreamParams()
		ex.RunGetUpstream(gp, up)
	},
}

var upstreamUpdateCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Update upstream",
	Long:  `Update upstream configuration object using config file and/or command line arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		LoadUpstreamParamsFromFile()
		ex := common.GetExecutor()
		gp := GetGlobalFlags()
		up := GetUpstreamParams()
		ex.RunUpdateUpstream(gp, up)
	},
}

func init() {
	createCmd.AddCommand(upstreamCmd)
	deleteCmd.AddCommand(upstreamDelCmd)
	describeCmd.AddCommand(upstreamDescribeCmd)
	getCmd.AddCommand(upstreamGetCmd)
	updateCmd.AddCommand(upstreamUpdateCmd)
	CreateNameParam(upstreamCmd, upstreamDelCmd, upstreamDescribeCmd, upstreamGetCmd, upstreamUpdateCmd)
	CreateTypeParam(upstreamCmd, upstreamUpdateCmd)
	CreateSpecParams(upstreamCmd, upstreamUpdateCmd)
}
