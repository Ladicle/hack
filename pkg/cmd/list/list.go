package list

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/util"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list contests",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := listDirLv2(config.BaseDir())
			if err != nil {
				return err
			}
			for _, c := range l {
				fmt.Println(c)
			}
			return nil
		},
		SilenceUsage: true,
	}
}

// listDirLv2 list visible directories up to 2 levels in the base
func listDirLv2(base string) ([]string, error) {
	var list []string
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
			list = append(list,
				filepath.Join(lv1.Name(), lv2.Name()))
		}
	}
	return list, nil
}
