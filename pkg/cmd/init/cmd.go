package init

import (
	"os"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/format"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

type initConfig struct {
	Number int
	*format.HackRobot
}

func NewCommand() *cobra.Command {
	ic := initConfig{HackRobot: format.NewHackRobot(os.Stdout)}

	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize workspace for the contest",
		Run:     ic.run,
	}
	cmd.Flags().IntVarP(&ic.Number, "number", "n", 4, "The number of quiz")
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
	case contest.HostCodeJam:
		cid := config.CurrentContestID()
		ic.Info("Prepare for the %q contest.\n", cid)
		n := 4
		if strings.Contains(strings.ToLower(cid), "round1") {
			n = 3
		}
		ic.Start("Creating %v quiz directories :package:", n)
		dirs := strings.Split(alphabet, "")[:n]
		if err := contest.MkQuizDir(dirs); err != nil {
			ic.Error()
			glog.Fatal(err)
		}
		ic.Success()

	case contest.HostFree:
		ic.Info("Prepare for the %q contest.\n", config.CurrentContestID())
		ic.Start("Creating %v quiz directories :package:", ic.Number)
		dirs := strings.Split(alphabet, "")[:ic.Number]
		if err := contest.MkQuizDir(dirs); err != nil {
			ic.Error()
			glog.Fatal(err)
		}
		ic.Success()
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
