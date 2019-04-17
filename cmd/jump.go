package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/util"
	"github.com/spf13/cobra"
)

func NewJumpCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "jump [quiz]",
		Aliases: []string{"j"},
		Short:   "Get current quiz directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Display the specified quiz path
			if len(args) >= 1 {
				fmt.Println(config.CurrentQuizPath(args[0]))
				return nil
			}
			wd, err := os.Executable()
			if err != nil {
				return err
			}
			cc := config.CurrentContestPath()
			s := strings.TrimPrefix(wd, cc)

			fs, err := ioutil.ReadDir(cc)
			if err != nil {
				return err
			}

			// In question directory
			if s != wd && s != "" {
				s = strings.TrimPrefix(s, "/")
				for i, f := range fs {
					if !util.IsVisibleDir(f) {
						continue
					}
					if f.Name() != s {
						continue
					}
					if i+1 == len(fs) {
						return fmt.Errorf("there is no more quiz")
					}
					fmt.Println(config.CurrentQuizPath(fs[i+1].Name()))
					return nil
				}
				return fmt.Errorf("%q is unexpected path", s)
			}

			// In other directory
			for _, f := range fs {
				if !util.IsVisibleDir(f) {
					continue
				}
				fmt.Println(config.CurrentQuizPath(f.Name()))
				return nil
			}
			return fmt.Errorf("%q is unexpected working directory", wd)
		},
		SilenceUsage: true,
	}
}
