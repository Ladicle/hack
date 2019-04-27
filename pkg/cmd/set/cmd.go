package set

import (
	"os"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/format"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

const pathFormat = "[host]/[contest_id]"

type setConfig struct{ *format.HackRobot }

func NewCommand() *cobra.Command {
	sc := setConfig{format.NewHackRobot(os.Stdout)}
	return &cobra.Command{
		Use:     "set " + pathFormat,
		Aliases: []string{"s"},
		Short:   "Switch contest current contest",
		Run: func(cmd *cobra.Command, args []string) {
			sc.run(args)
		},
		SilenceUsage: true,
	}
}

func (sc *setConfig) run(args []string) {
	if len(args) != 1 {
		sc.Fatal("Please set required argument %q.", pathFormat)
	}

	c := args[0]
	if len(strings.Split(c, "/")) != 2 {
		sc.Fatal("Oops, %q is invalid path. You can set %q.", c, pathFormat)
	}

	config.SetCurrent(c)
	sc.Info("OK! I set %q for the next contest", c)

	path := config.CurrentContestPath()
	if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
		glog.V(8).Infof("%v is already exists: err=%v", path, err)
		return
	}
	sc.Info("And... it's the first time. Good luck!\n")

	sc.Start("Creating contest directory :package:")
	defer sc.End()
	if err := contest.MkCurrentContestDir(); err != nil {
		sc.Fatal("Sorry, failed to create %q directory", path)
		glog.Error(err)
	}
}
