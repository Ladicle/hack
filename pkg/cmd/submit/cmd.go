package submit

import (
	"os"
	"time"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/format"
	"github.com/Ladicle/hack/pkg/util"
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
	quizID, err := contest.CurrentQuizID(wd)
	if err != nil {
		cfg.Fatal("Not in the contest directory.")
	}

	sourceFile, err := util.GetProgName()
	if err != nil {
		glog.Fatal(err)
	}
	if sourceFile == "" {
		cfg.Fatal("Not found the %v program in this directory.", "main.[go|cpp]")
	}

	s := NewSolution(quizID, sourceFile, cfg.HackRobot)
	if err := s.Test(cfg.TimeLimit); err != nil {
		glog.Fatal(err)
	}

	cfg.Printfln("")
	if config.CurrentHost() == contest.HostAtCoder {
		if err := s.Submit(); err != nil {
			glog.Fatal(err)
		}
		if err := s.Open(); err != nil {
			glog.Fatal(err)
		}
	}
	if err := s.Copy(); err != nil {
		glog.Fatal(err)
	}
}
