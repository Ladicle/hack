package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
				fmt.Println(config.SetCurrentQuizPath(args[0]))
				return nil
			}

			// In question directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			cc := config.CurrentContestPath()

			fs, err := ioutil.ReadDir(cc)
			if err != nil {
				return err
			}

			if isInQuizDir(wd, cc) {
				quiz := filepath.Base(wd)
				for i, f := range fs {
					if !util.IsVisibleDir(f) {
						continue
					}
					if f.Name() != quiz {
						continue
					}
					if i+1 == len(fs) {
						// There is no more quiz
						fmt.Println(config.SetCurrentQuizPath(fs[i].Name()))
						return nil
					}
					fmt.Println(config.SetCurrentQuizPath(fs[i+1].Name()))
					return nil
				}
				return fmt.Errorf("%q is unexpected path", quiz)
			}

			// In other directory
			for _, f := range fs {
				if !util.IsVisibleDir(f) {
					continue
				}
				fmt.Println(config.SetCurrentQuizPath(f.Name()))
				return nil
			}
			return fmt.Errorf("%q is unexpected working directory", wd)
		},
		SilenceUsage: true,
	}
}

func isInQuizDir(workingDir, contestPath string) bool {
	s := strings.TrimPrefix(workingDir, contestPath)
	return s != workingDir && s != ""
}
