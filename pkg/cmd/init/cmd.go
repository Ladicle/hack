package init

import (
	"os"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/format"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

type initConfig struct {
	*format.HackRobot
}

func NewCommand() *cobra.Command {
	ic := initConfig{format.NewHackRobot(os.Stdout)}

	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize workspace for the contest",
		Run:     ic.run,
	}
	return cmd
}

func (ic *initConfig) run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		glog.Fatal(err)
	}
	ch := config.CurrentHost()
	cq, err := contest.CurrentQuizID(wd)
	if err != nil {
		if wd == config.CurrentContestPath() {
			ic.initContest(ch)
			return
		}
		ic.Fatal("Not in the contest directory.")
	}
	ic.initQuiz(ch, cq)
}

func (ic *initConfig) initContest(host string) {
	switch host {
	case contest.HostAtCoder:
		at := newAtCoderInitializer(ic.HackRobot, config.CurrentContestID())
		if err := at.initAtCoderContest(); err != nil {
			glog.Fatal(err)
		}
	default:
		ic.Fatal("Sorry, %q is not supported.", host)
	}
}

func (ic *initConfig) initQuiz(host, quiz string) {
	switch host {
	case contest.HostAtCoder:
		at := newAtCoderInitializer(ic.HackRobot, config.CurrentContestID())
		if err := at.createSamples(quiz); err != nil {
			glog.Fatal(err)
		}
	default:
		ic.Fatal("Sorry, %q is not supported.", host)
	}
}
