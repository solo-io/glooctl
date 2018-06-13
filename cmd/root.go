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
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap/flags"
	"github.com/solo-io/glooctl/cmd/install"
	"github.com/solo-io/glooctl/cmd/route"
	"github.com/solo-io/glooctl/cmd/secret"
	"github.com/solo-io/glooctl/cmd/upstream"
	"github.com/solo-io/glooctl/cmd/virtualservice"
	"github.com/solo-io/glooctl/pkg/config"
	"github.com/spf13/cobra"
)

var (
	// Functions like __glooctl_get_* etc are custom functions in
	// bash script to get completion parameters. It uses `compgen` to
	// setup COMPREPLY variable.
	// compgen takes a space separate list of options. We use glooctl
	// itself to genertae these options. We use -o template to format
	// the glooctl output.
	//
	// for additional documentation on how auto completion works
	// please see https://github.com/spf13/cobra/blob/master/bash_completions.md
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

__glooctl_get_virtualservices()
{
	local glooctl_out
	if glooctl_out=$(glooctl virtualservice get -o template --template="{{range .}}{{.Name}} {{end}}" 2>/dev/null); then
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
		glooctl_virtualservice_edit | glooctl_virtualservice_delete | glooctl_virtualservice_get)
			__glooctl_get_virtualservices
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

	config.LoadConfig(opts) // load saved configurations

	app.SuggestionsMinimumDistance = 1
	app.AddCommand(
		configcmd.Cmd(opts),
		upstream.Cmd(opts),
		functionCmd(opts),
		virtualservice.Cmd(opts),
		route.Cmd(opts),
		secret.Cmd(opts),
		install.Cmd(),
		registerCmd(opts),
		completionCmd())

	return app
}
