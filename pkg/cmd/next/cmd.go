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

var reversed bool

func NewCommand(f *config.File, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "next",
		Aliases:      []string{"n"},
		Short:        "Print next quiz directory",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	cmd.Flags().BoolVarP(&reversed, "reversed", "r", false, "Prints previous problem")

	return cmd
}

func run() error {
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
	for i, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if next {
			fmt.Print(filepath.Join(contestAbsDir, entry.Name()))
			return nil
		}
		if entry.Name() == curTaskDir {
			if reversed {
				if i-1 >= 0 {
					fmt.Print(filepath.Join(contestAbsDir, entries[i-1].Name()))
					return nil
				}
				return fmt.Errorf("%s is first task", curTaskDir)
			}
			next = true
		}
	}
	return fmt.Errorf("%s is last task", curTaskDir)
}
