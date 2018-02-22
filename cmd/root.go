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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/solo-io/glooctl/cmd/config"
	"github.com/solo-io/glooctl/cmd/route"
	"github.com/solo-io/glooctl/cmd/upstream"
	"github.com/solo-io/glooctl/cmd/vhost"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	glooFolder string
	kubeConfig string
	namespace  string
	syncPeriod int
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
	// TODO(ashish) Enable Viper
	//cobra.OnInitialize(initConfig)

	// global flags
	flags := rootCmd.PersistentFlags()
	//flags.StringVar(&cfgFile, "glooconfig", "", "config file (default is $HOME/.glooctl.yaml)")
	flags.StringVar(&glooFolder, "gloo-folder", "", "storage folder")
	flags.StringVar(&kubeConfig, "kubeconfig", "", "kubeconfig (defaults to ~/.kube/config)")
	flags.StringVarP(&namespace, "namespace", "n", "", "namespace for resources")
	flags.IntVarP(&syncPeriod, "sync-period", "s", 60, "sync period (seconds) for resources")

	rootCmd.SuggestionsMinimumDistance = 1
	rootCmd.AddCommand(
		upstream.UpstreamCmd(),
		vhost.VHostCmd(),
		route.RouteCmd(),
		config.ConfigCmd(),
		registerCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".glooctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".glooctl")
	}

	viper.SetEnvPrefix("gloo")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Unable to read config: ", err)
		os.Exit(1)
	}
}
