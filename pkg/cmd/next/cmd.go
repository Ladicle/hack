package next

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand(f *config.File, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:          "next",
		Aliases:      []string{"n"},
		Short:        "Print next quiz directory",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			var (
				curTaskDir    = filepath.Base(wd)
				contestAbsDir = filepath.Dir(wd)
			)
			entries, err := os.ReadDir(contestAbsDir)
			if err != nil {
				return err
			}

			sort.Slice(entries, func(i, j int) bool {
				return entries[i].Name() < entries[j].Name()
			})

			var next bool
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				if next {
					fmt.Print(filepath.Join(contestAbsDir, entry.Name()))
					return nil
				}
				if entry.Name() == curTaskDir {
					next = true
				}
			}
			return fmt.Errorf("%s is last task", curTaskDir)
		},
	}
}
