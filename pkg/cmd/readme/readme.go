package readme

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/readme"
)

func NewCommand(f *config.File) *cobra.Command {
	return &cobra.Command{
		Use:          "readme <CONTEST_ID> <TASK_ID>",
		Aliases:      []string{"r"},
		Short:        "Generate README.org in quiz directory",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if got, want := len(args), 2; got != want {
				return fmt.Errorf("invalid number of arguments: got=%v, want=%v", got, want)
			}

			var (
				contestID = args[0]
				taskID    = args[1]
			)
			data, err := readme.GenerateReadme(contestID, taskID, time.Now())
			if err != nil {
				return err
			}

			dst, err := contest.GetDirFromID(f.BaseDir, contestID, taskID)
			if err != nil {
				return err
			}
			return os.WriteFile(dst, data, 0666)
		},
	}
}
