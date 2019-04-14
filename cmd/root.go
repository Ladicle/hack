package cmd

import (
	"github.com/spf13/cobra"
)

var (
	gitRepo string
	version string
)

var rootCmd = &cobra.Command{
	Use:   "hack",
	Short: "Hack assists your programming contest.",
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(sampleCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
