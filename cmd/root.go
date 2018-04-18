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
	"github.com/solo-io/glooctl/cmd/route"
	"github.com/solo-io/glooctl/cmd/secret"
	"github.com/solo-io/glooctl/cmd/upstream"
	"github.com/solo-io/glooctl/cmd/vhost"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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

func App(version string) *cobra.Command {
	app := &cobra.Command{
		Use:   "glooctl",
		Short: "manage resources in the Gloo Universe",
		Long: `glooctl configures resources use by Gloo server.
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

	opts.ConfigStorageOptions.Type = "kube"
	opts.SecretStorageOptions.Type = "kube"
	opts.FileStorageOptions.Type = "kube"
	opts.KubeOptions.KubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")

	app.SuggestionsMinimumDistance = 1
	app.AddCommand(
		upstream.UpstreamCmd(opts),
		functionCmd(opts),
		vhost.VHostCmd(opts),
		route.RouteCmd(opts),
		secret.SecretCmd(opts),
		registerCmd(opts),
		completionCmd())

	return app
}
