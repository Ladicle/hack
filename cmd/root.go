package cmd

import (
	"fmt"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

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
		fmt.Sprintf("config file (default %v)", config.DefaultCfg()))

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewSampleCmd())
	rootCmd.AddCommand(NewCopyCmd())
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewSwitchCmd())
	rootCmd.AddCommand(NewListCmd())
	rootCmd.AddCommand(NewJumpCmd())
	rootCmd.AddCommand(NewTestCmd())
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}

func initConfig() {
	config.Load(cfgFile)
}

func Execute() error {
	return rootCmd.Execute()
}
