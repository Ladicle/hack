package submit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/lang"
)

type Options struct {
	Timeout time.Duration
	DryRun  bool
}

func NewCommand(f *config.File, out io.Writer) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:          "submit",
		Aliases:      []string{"sub", "s"},
		Short:        "Submit the solution",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run(f, out)
		},
	}

	cmd.Flags().DurationVar(&opts.Timeout, "timeout", 2*time.Second, "Set timeout duration.")
	cmd.Flags().BoolVarP(&opts.DryRun, "dry-run", "d", false, "Only run tests and do not submit program.")

	return cmd
}

func (o *Options) Run(f *config.File, out io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Test program with all samples
	prog, err := findMainFile(wd)
	if err != nil {
		return err
	}
	tester, err := lang.GetTester(prog, o.Timeout)
	if err != nil {
		return err
	}
	sampleNum, err := cntSamples(wd)
	if err != nil {
		return err
	}
	var cntErr int
	for sampleID := 1; sampleID <= sampleNum; sampleID++ {
		err = tester.Run(sampleID)
		if err == nil {
			fmt.Fprintf(out, "[AC] Sample #%d\n", sampleID)
			continue
		}
		var langErr lang.Error
		if ok := errors.As(err, &langErr); !ok {
			return err
		}
		fmt.Fprintln(out, langErr)
		cntErr++
	}
	if cntErr > 0 {
		return fmt.Errorf("Fail %d samples", cntErr)
	}

	// Submit program
	var (
		contestID = getContestID(wd)
		taskID    = getTaskID(wd)
	)
	if o.DryRun {
		fmt.Fprintf(out, "Submit %v (dry-run)\n", taskID)
		return nil
	}
	at, err := contest.NewAtCoder(contestID)
	if err != nil {
		return err
	}
	if err := at.Login(f.AtCoder.User, f.AtCoder.Pass); err != nil {
		return err
	}
	if err := at.Submit(taskID, prog); err != nil {
		return err
	}
	fmt.Fprintf(out, "Submit %v\n", taskID)
	return nil
}

// findMainFile finds a file that has 'main' as a name prefix in the current directory.
func findMainFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(filepath.Base(entry.Name()), "main") {
			return entry.Name(), nil
		}
	}
	return "", errors.New("not found 'main' program in the current directory")
}

// cntSamples counts the number of sample input files which has a '.in' extension.
func cntSamples(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return -1, err
	}
	var cnt int
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".in" {
			cnt++
		}
	}
	return cnt, nil
}

// getContestID returns the parent directory name as the contest ID.
func getContestID(dir string) string {
	curBase := filepath.Base(dir)
	return filepath.Base(dir[:len(dir)-len(curBase)])
}

// getTaskID returns the specified directory name as the task ID.
func getTaskID(dir string) string {
	return filepath.Base(dir)
}
