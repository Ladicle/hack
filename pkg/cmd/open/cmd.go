package open

import (
	"os"

	"github.com/Ladicle/hack/pkg/contest"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "open",
		Aliases:      []string{"o"},
		Short:        "Open current task page",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			var (
				contestID = contest.GetContestID(wd)
				taskID    = contest.GetTaskID(wd)
			)
			return browser.OpenURL(contest.GetTaskURL(contestID, taskID))
		},
	}
}
