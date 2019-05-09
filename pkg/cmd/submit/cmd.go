package submit

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/format"
	"github.com/Ladicle/hack/pkg/util"
	"github.com/atotto/clipboard"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

type submitConig struct {
	TimeLimit time.Duration

	*format.HackRobot
}

func NewCommand() *cobra.Command {
	cfg := submitConig{HackRobot: format.NewHackRobot(os.Stdout)}
	cmd := &cobra.Command{
		Use:     "submit",
		Aliases: []string{"sub"},
		Short:   "Submit the solution",
		Run:     cfg.run,
	}
	cmd.Flags().DurationVarP(&cfg.TimeLimit, "time-limit", "t", 2*time.Second,
		"Set execution time-limit")
	return cmd
}

func (cfg *submitConig) run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		glog.Fatal(err)
	}
	if _, err := contest.CurrentQuizID(wd); err != nil {
		cfg.Fatal("Not in the contest directory.")
	}

	if err := runTest(cfg.TimeLimit, cfg.HackRobot); err != nil {
		glog.Fatal(err)
	}

	if config.CurrentHost() == contest.HostAtCoder {
		// submit
	} else {
		if err := copyProgram(cfg.HackRobot); err != nil {
			glog.Fatal(err)
		}
	}
}

func copyProgram(hr *format.HackRobot) error {
	fname, err := util.GetProgName()
	if err != nil {
		return err
	}
	if fname == "" {
		hr.Fatal("Not found the %v program in this directory.", "main.[go|cpp]")
	}

	code, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(string(code)); err != nil {
		return err
	}
	hr.Info("\nCopy %q to the clipboard.", fname)
	return nil
}
