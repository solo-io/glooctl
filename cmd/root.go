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
	"github.com/solo-io/glooctl/cmd/route"
	"github.com/solo-io/glooctl/cmd/secret"
	"github.com/solo-io/glooctl/cmd/upstream"
	"github.com/solo-io/glooctl/cmd/vhost"
	"github.com/spf13/cobra"
)

var (
	cfgFile        string
	resourceFolder string
	secretFolder   string
	kubeConfig     string
	namespace      string
	syncPeriod     int
)

func App(version string) *cobra.Command {
	app := &cobra.Command{
		Use:   "glooctl",
		Short: "manage resources in the Gloo Universe",
		Long: `glooctl configures resources use by Gloo server.
	Find more information at https://gloo.solo.io`,
		Version: version,
	}
	// global flags
	flags := app.PersistentFlags()
	flags.StringVar(&kubeConfig, "kubeconfig", "", "kubeconfig (defaults to ~/.kube/config)")
	flags.StringVarP(&namespace, "namespace", "n", "gloo-system", "namespace for resources")
	flags.StringVar(&resourceFolder, "gloo-config-dir", "", "if set, glooctl will use file-based storage. use this if gloo is running locally, "+
		"e.g. using docker with volumes mounted for config storage.")
	flags.StringVar(&secretFolder, "secret-dir", "", "if set, glooctl will use file-based stroage. use this if gloo is running locally")
	flags.IntVarP(&syncPeriod, "sync-period", "s", 60, "sync period (seconds) for resources")

	app.SuggestionsMinimumDistance = 1
	app.AddCommand(
		upstream.UpstreamCmd(),
		vhost.VHostCmd(),
		route.RouteCmd(),
		secret.SecretCmd(),
		registerCmd())

	return app
}
