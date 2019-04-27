package get

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/format"
	"github.com/Ladicle/hack/pkg/util"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "get contests",
		RunE: func(cmd *cobra.Command, args []string) error {
			hb := format.NewHackRobot(os.Stdout)
			cc := config.CurrentContest()

			hb.Info("Current contest is %q", cc)

			l, err := getDirLv2(config.BaseDir())
			if err != nil {
				return err
			}

			hb.Info("There are %v other contests like this:\n", len(l)-1)
			for _, c := range l {
				if c == cc {
					continue
				}
				fmt.Println(" - " + c)
			}
			return nil
		},
		SilenceUsage: true,
	}
}

// getDirLv2 gets visible directories up to 2 levels in the base
func getDirLv2(base string) ([]string, error) {
	var dirs []string
	lv1fs, err := ioutil.ReadDir(base)
	if err != nil {
		return nil, err
	}

	for _, lv1 := range lv1fs {
		if !util.IsVisibleDir(lv1) {
			continue
		}

		lv2fs, err := ioutil.ReadDir(filepath.Join(base, lv1.Name()))
		if err != nil {
			return nil, err
		}

		for _, lv2 := range lv2fs {
			if !util.IsVisibleDir(lv2) {
				continue
			}
			dirs = append(dirs,
				filepath.Join(lv1.Name(), lv2.Name()))
		}
	}
	return dirs, nil
}
