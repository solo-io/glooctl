// Copyright © 2018 Solo.io <anton.stadnikov@solo.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/solo-io/gluectl/platform/common"
	"github.com/spf13/cobra"
)

// upstreamUpdateCmd represents the upstreamUpdate command
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
	updateCmd.AddCommand(upstreamUpdateCmd)
	CreateUpstreamParams(upstreamUpdateCmd, true)
}
