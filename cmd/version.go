package cmd

import (
	"fmt"
	"git-projects/internal/name"
	"git-projects/internal/version"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd returns a new version command.
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   fmt.Sprintf("Show the %s version information", name.ApplicationName),
	Example: fmt.Sprintf("%s version", name.ApplicationName),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf(`
=======================================================================
=                           %s                              =
=======================================================================
version     : %s
commit      : %s
date        : %s
go version  : %s
go compiler : %s
platform    : %s/%s
`, name.ApplicationName, version.Version, version.Commit, version.Date, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
		return nil
	},
}

// Initialize command
func init() {
	rootCmd.AddCommand(versionCmd)
}
