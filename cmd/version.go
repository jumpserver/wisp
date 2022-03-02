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
		fmt.Printf("Wisp Version: %s\n", Version)
		fmt.Printf("Go Version: %s\n", GoVersion)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build UTC Time: %s\n", BuildTime)
	},
}
