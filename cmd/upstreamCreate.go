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
	"fmt"
	"log"

	"github.com/solo-io/gluectl/platform"
	"github.com/spf13/cobra"
)

// upstreamCmd represents the upstream command
var upstreamCmd = &cobra.Command{
	Use:   "upstream",
	Short: "Create upstream",
	Long: `Create upstream configuration object using config file or command line arguments
	Example: gluectl create upstream --file MyObj.yaml
	Example: gluectl create upstream --name "Name" --type "k8s" --spec "spec.service='svc';spec.port=8080"`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upstream called")
		ex := platform.GetExecutor()

		file, namespace, wait := GetGlobalFlags(cmd)
		if file != "" {
			ex.RunCreateUpstreamFromFile(file, namespace, wait)
		} else {
			name, err := cmd.LocalFlags().GetString("name")
			if err != nil {
				log.Fatal("Invalid value of the 'name' flag", err)
			}
			t, err := cmd.LocalFlags().GetString("type")
			if err != nil {
				log.Fatal("Invalid value of the 'type' flag", err)
			}
			spec, err := cmd.LocalFlags().GetString("spec")
			if err != nil {
				log.Fatal("Invalid value of the 'spec' flag", err)
			}
			ex.RunCreateUpstream(name, namespace, t, spec, wait)
		}
	},
}

func init() {
	createCmd.AddCommand(upstreamCmd)
	upstreamCmd.PersistentFlags().String("name", "", "upstream name")
	upstreamCmd.PersistentFlags().String("type", "", "upstream type")
	upstreamCmd.PersistentFlags().String("spec", "", "key-value pairs separated by ';', ex. --spec \"k1=v1;k2=v2\"")

}
