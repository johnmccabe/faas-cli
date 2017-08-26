// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"github.com/spf13/cobra"
)

const defaultNetwork = "func_functions"

// Flags that are to be added to all commands.
var (
	yamlFile string
)

// Flags that are to be added to subset of commands.
var (
	fprocess     string
	functionName string
	gateway      string
	handler      string
	image        string
	language     string
)

func init() {
	faasCmd.PersistentFlags().StringVarP(&yamlFile, "yaml", "f", "", "Path to YAML file describing function(s)")

	// Set Bash completion options
	validYAMLFilenames := []string{"yaml", "yml"}
	_ = faasCmd.PersistentFlags().SetAnnotation("yaml", cobra.BashCompFilenameExt, validYAMLFilenames)
}

// Execute TODO
func Execute() {
	faasCmd.SilenceUsage = true
	faasCmd.Execute()
}

// faasCmd is the FaaS CLI root command and mimics the legacy client behaviour
// Every other command attached to FaasCmd is a child command to it.
var faasCmd = &cobra.Command{
	Use:   "faas-cli",
	Short: "build and deploy FaaS functions",
	Long: `TODO Add a full description including a link to
the commands documentation - https://github.com/alexellis/faas`,
	Run: runFaas,
}

// runFaas TODO
func runFaas(cmd *cobra.Command, args []string) {
	cmd.Help()
}
