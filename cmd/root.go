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
	"os"

	"github.com/solo-io/glooctl/cmd/route"
	"github.com/solo-io/glooctl/cmd/upstream"
	"github.com/solo-io/glooctl/cmd/vhost"
	"github.com/spf13/cobra"
)

var (
	cfgFile        string
	resourceFolder string
	kubeConfig     string
	namespace      string
	syncPeriod     int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "glooctl",
	Short: "manage resources in the Gloo Universe",
	Long: `glooctl configures upstreams and virtual hosts to be used by Gloo server
	Find more information at https://github.com/solo-io/gloo`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// global flags
	flags := rootCmd.PersistentFlags()
	flags.StringVar(&kubeConfig, "kubeconfig", "", "kubeconfig (defaults to ~/.kube/config)")
	flags.StringVarP(&namespace, "namespace", "n", "gloo-system", "namespace for resources")
	flags.StringVar(&resourceFolder, "resource-folder", "", "folder for storing resources when using file based store")
	flags.IntVarP(&syncPeriod, "sync-period", "s", 60, "sync period (seconds) for resources")

	rootCmd.SuggestionsMinimumDistance = 1
	rootCmd.AddCommand(
		upstream.UpstreamCmd(),
		vhost.VHostCmd(),
		route.RouteCmd(),
		registerCmd())
}
