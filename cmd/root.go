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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"

	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/flags"
	"github.com/solo-io/glooctl/cmd/route"
	"github.com/solo-io/glooctl/cmd/secret"
	"github.com/solo-io/glooctl/cmd/upstream"
	"github.com/solo-io/glooctl/cmd/vhost"
	"github.com/solo-io/glooctl/pkg/util"
	"github.com/spf13/cobra"
)

const (
	defaultStorage = "kube"
)

var (
	bashCompletion = `
__glooctl_route_http_methods()
{
	COMPREPLY=( $( compgen -W "GET PUT POST PATCH DELETE HEAD OPTIONS" -- "$cur" ) )
}

__glooctl_get_upstreams()
{
	local glooctl_out
	if glooctl_out=$(glooctl upstream get -o template --template="{{range .}}{{.Name}} {{end}}" 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${glooctl_out}" -- "$cur" ) )
	fi
}

__glooctl_get_functions()
{
	local glooctl_out
	if glooctl_out=$(glooctl function get -o template --template="{{range .}}{{.Function.Name}} {{end}}" 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${glooctl_out}" -- "$cur" ) )
	fi
}

__glooctl_get_virtualhosts()
{
	local glooctl_out
	if glooctl_out=$(glooctl virtualhost get -o template --template="{{range .}}{{.Name}} {{end}}" 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${glooctl_out}" -- "$cur" ) )
	fi
}

__custom_func()
{
	case ${last_command} in
		glooctl_upstream_edit | glooctl_upstream_delete | glooctl_upstream_get)
			__glooctl_get_upstreams
			return
			;;
		glooctl_virtualhost_edit | glooctl_virtualhost_delete | glooctl_virtualhost_get)
			__glooctl_get_virtualhosts
			return
			;;
		*)
			;;
	esac
}
	`
)

// App returns the command representing the CLI
func App(version string) *cobra.Command {
	app := &cobra.Command{
		Use:   "glooctl",
		Short: "manage resources in the Gloo Universe",
		Long: `glooctl configures resources used by Gloo server.
	Find more information at https://gloo.solo.io`,
		Version:                version,
		BashCompletionFunction: bashCompletion,
	}

	opts := &bootstrap.Options{}
	flags.AddConfigStorageOptionFlags(app, opts)
	flags.AddSecretStorageOptionFlags(app, opts)
	flags.AddFileStorageOptionFlags(app, opts)

	flags.AddFileFlags(app, opts)
	flags.AddKubernetesFlags(app, opts)
	flags.AddConsulFlags(app, opts)
	flags.AddVaultFlags(app, opts)

	loadConfig(opts) // load saved configurations

	app.SuggestionsMinimumDistance = 1
	app.AddCommand(
		upstream.UpstreamCmd(opts),
		functionCmd(opts),
		vhost.VHostCmd(opts),
		route.Cmd(opts),
		secret.Cmd(opts),
		registerCmd(opts),
		completionCmd())

	return app
}

// loadConfig loads saved configuration if any
// if not sets default configuration and also saves it
func loadConfig(opts *bootstrap.Options) {
	configDir, err := util.ConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to get config directory:", err)
		defaultConfig(opts)
		return
	}
	configFile := filepath.Join(configDir, "config.yaml")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		defaultConfig(opts)
		if os.IsNotExist(err) {
			saveConfig(opts, configFile)
		} else {
			fmt.Fprintln(os.Stderr, "Error reading configuration file:", err)
		}
		return
	}
	if err := yaml.Unmarshal(data, opts); err != nil {
		defaultConfig(opts)
		fmt.Fprintln(os.Stderr, "Unable to parse configuration file:", err)
	}
}

func saveConfig(opts *bootstrap.Options, configFile string) {
	b, err := yaml.Marshal(opts)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(configFile, b, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to save configuration file", configFile)
	}
}

func defaultConfig(opts *bootstrap.Options) {
	opts.ConfigStorageOptions.Type = defaultStorage
	opts.SecretStorageOptions.Type = defaultStorage
	opts.FileStorageOptions.Type = defaultStorage
	opts.KubeOptions.KubeConfig = filepath.Join(util.HomeDir(), ".kube", "config")
}
