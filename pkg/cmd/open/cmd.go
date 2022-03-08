package open

import (
	"fmt"
	"os"

	"github.com/Ladicle/hack/pkg/contest"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var browse bool

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
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
				addr      = contest.GetTaskURL(contestID, taskID)
			)
			if browse {
				return browser.OpenURL(addr)
			}
			fmt.Print(addr)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&browse, "browse", "b", true, "Open task in browser. If false, the URL will be output.")
	return cmd
}
