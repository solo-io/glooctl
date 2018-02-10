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
		InteractiveModeUpstream("create")
		common.GetExecutor().RunCreateUpstream(GetGlobalFlags(), GetUpstreamParams())
	},
}

var upstreamDelCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Delete upstream",
	Long:  `Delete Upstream by name`,
	Run: func(cmd *cobra.Command, args []string) {
		LoadUpstreamParamsFromFile()
		InteractiveModeUpstream("delete")
		common.GetExecutor().RunDeleteUpstream(GetGlobalFlags(), GetUpstreamParams())
	},
}

var upstreamDescribeCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Describe upstream(s)",
	Long:  `Describe upstream (by name) or all upstreams in the namespace`,

	Run: func(cmd *cobra.Command, args []string) {
		InteractiveModeUpstream("describe")
		common.GetExecutor().RunDescribeUpstream(GetGlobalFlags(), GetUpstreamParams())
	},
}

var upstreamGetCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Get upstream(s)",
	Long:  `Get upstream (by name) or all upstreams in the namespace`,

	Run: func(cmd *cobra.Command, args []string) {
		InteractiveModeUpstream("get")
		common.GetExecutor().RunGetUpstream(GetGlobalFlags(), GetUpstreamParams())
	},
}

var upstreamUpdateCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Update upstream",
	Long:  `Update upstream configuration object using config file and/or command line arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		LoadUpstreamParamsFromFile()
		InteractiveModeUpstream("update")
		common.GetExecutor().RunUpdateUpstream(GetGlobalFlags(), GetUpstreamParams())
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
