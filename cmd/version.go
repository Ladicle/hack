package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show this command version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v -- %v\n", gitRepo, version)
	},
}
