package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"

	getcmd "github.com/Ladicle/hack/pkg/cmd/get"
	initcmd "github.com/Ladicle/hack/pkg/cmd/init"
	jumpcmd "github.com/Ladicle/hack/pkg/cmd/jump"
	setcmd "github.com/Ladicle/hack/pkg/cmd/set"
	smtcmd "github.com/Ladicle/hack/pkg/cmd/submit"
	"github.com/Ladicle/hack/pkg/config"
)

const welcomeMsg = `Welcome to the Hack!:tada:

:robot:< I'm a hack bot!
   < Before start hacking, please answer some questions.
`

var (
	gitRepo string
	version string

	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "hack",
	Short: "Hack assists your programming contest.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		flag.Parse()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		glog.V(4).Info("Saving configuration...")
		err := config.Save()
		glog.V(4).Info("Saved configuration")
		return err
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		fmt.Sprintf("config file (default ~/%v)", config.DefaultConfig))

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(getcmd.NewCommand())
	rootCmd.AddCommand(jumpcmd.NewCommand())
	rootCmd.AddCommand(setcmd.NewCommand())
	rootCmd.AddCommand(initcmd.NewCommand())
	rootCmd.AddCommand(smtcmd.NewCommand())

	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}

func initConfig() {
	err := config.Load(cfgFile)
	if err == nil {
		return
	}
	if !os.IsNotExist(err) {
		glog.Fatal(err)
	}

	// set defaults at first time
	emoji.Println(welcomeMsg)

	dir, err := askBaseDir()
	if err != nil {
		glog.Fatal(err)
	}
	config.SetBaseDir(dir)

	var ans string
	fmt.Printf("# Do you have AtCoder account? (y/n)\n-> ")
	fmt.Scanln(&ans)
	if ans == "y" {
		ac, err := initAtCoder()
		if err != nil {
			glog.Fatal(err)
		}
		config.SetAtCoderAccount(ac)
	} else {
		fmt.Printf("OK! I'll skip this.\n\n")
	}
}

func Execute() error {
	return rootCmd.Execute()
}
