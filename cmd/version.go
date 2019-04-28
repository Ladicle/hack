package cmd

import (
	"os"

	"github.com/Ladicle/hack/pkg/format"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Show this command version",
		Run: func(cmd *cobra.Command, args []string) {
			hb := format.NewHackRobot(os.Stdout)
			hb.Info("Hi, I'm %v robot! Current version is %v.",
				gitRepo, version)
		},
	}
}
