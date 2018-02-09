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

// upstreamGetCmd represents the upstreamGet command
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

func init() {
	getCmd.AddCommand(upstreamGetCmd)
	CreateUpstreamParams(upstreamGetCmd, false)

}
