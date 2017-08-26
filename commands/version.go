package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GitCommit injected at build-time
var GitCommit string

func init() {
	faasCmd.AddCommand(versionCmd)
}

// versionCmd displays version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the clients version information",
	Long: fmt.Sprintf(`The version command returns the current clients version
information.

This currently consists of the GitSHA from which the client was built.

See https://github.com/alexellis/faas-cli/tree/%s for the build source.`, GitCommit),
	Run: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Git Commit: %s\n", GitCommit)
	return
}
