/*
	Copyright (c) 2020 Docker Inc.

	Permission is hereby granted, free of charge, to any person
	obtaining a copy of this software and associated documentation
	files (the "Software"), to deal in the Software without
	restriction, including without limitation the rights to use, copy,
	modify, merge, publish, distribute, sublicense, and/or sell copies
	of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be
	included in all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
	EXPRESS OR IMPLIED,
	INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
	IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
	HOLDERS BE LIABLE FOR ANY CLAIM,
	DAMAGES OR OTHER LIABILITY,
	WHETHER IN AN ACTION OF CONTRACT,
	TORT OR OTHERWISE,
	ARISING FROM, OUT OF OR IN CONNECTION WITH
	THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package context

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/docker/api/cli/dockerclassic"
	"github.com/docker/api/context/store"
)

type descriptionCreateOpts struct {
	description string
}

func createCommand() *cobra.Command {
	const longHelp = `Create a new context

Create docker engine context: 
$ docker context create CONTEXT [flags]

Create Azure Container Instances context:
$ docker context create aci CONTEXT [flags]
(see docker context create aci --help)

Docker endpoint config:

NAME                DESCRIPTION
from                Copy named context's Docker endpoint configuration
host                Docker endpoint on which to connect
ca                  Trust certs signed only by this CA
cert                Path to TLS certificate file
key                 Path to TLS key file
skip-tls-verify     Skip TLS certificate validation

Kubernetes endpoint config:

NAME                 DESCRIPTION
from                 Copy named context's Kubernetes endpoint configuration
config-file          Path to a Kubernetes config file
context-override     Overrides the context set in the kubernetes config file
namespace-override   Overrides the namespace set in the kubernetes config file

Example:

$ docker context create my-context --description "some description" --docker "host=tcp://myserver:2376,ca=~/ca-file,cert=~/cert-file,key=~/key-file"`

	cmd := &cobra.Command{
		Use:   "create CONTEXT",
		Short: "Create new context",
		RunE: func(cmd *cobra.Command, args []string) error {
			return dockerclassic.ExecCmd(cmd)
		},
		Long: longHelp,
	}

	cmd.AddCommand(
		createAciCommand(),
		createLocalCommand(),
		createExampleCommand(),
	)

	flags := cmd.Flags()
	flags.String("description", "", "Description of the context")
	flags.String(
		"default-stack-orchestrator", "",
		"Default orchestrator for stack operations to use with this context (swarm|kubernetes|all)")
	flags.StringToString("docker", nil, "set the docker endpoint")
	flags.StringToString("kubernetes", nil, "set the kubernetes endpoint")
	flags.String("from", "", "create context from a named context")

	return cmd
}

func createLocalCommand() *cobra.Command {
	var opts descriptionCreateOpts
	cmd := &cobra.Command{
		Use:    "local CONTEXT",
		Short:  "Create a context for accessing local engine",
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createDockerContext(cmd.Context(), args[0], store.LocalContextType, opts.description, store.LocalContext{})
		},
	}
	addDescriptionFlag(cmd, &opts.description)
	return cmd
}

func createExampleCommand() *cobra.Command {
	var opts descriptionCreateOpts
	cmd := &cobra.Command{
		Use:    "example CONTEXT",
		Short:  "Create a test context returning fixed output",
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createDockerContext(cmd.Context(), args[0], store.ExampleContextType, opts.description, store.ExampleContext{})
		},
	}

	addDescriptionFlag(cmd, &opts.description)
	return cmd
}

func createDockerContext(ctx context.Context, name string, contextType string, description string, data interface{}) error {
	s := store.ContextStore(ctx)
	result := s.Create(
		name,
		contextType,
		description,
		data,
	)
	return result
}

func addDescriptionFlag(cmd *cobra.Command, descriptionOpt *string) {
	cmd.Flags().StringVar(descriptionOpt, "description", "", "Description of the context")
}
