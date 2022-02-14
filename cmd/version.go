package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "Unknown"
	GoVersion = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of wisp",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Wisp version: %s\n", Version)
		fmt.Printf("Go version: %s\n", GoVersion)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Built on: %s\n", BuildTime)
	},
}
