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
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewSampleCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
