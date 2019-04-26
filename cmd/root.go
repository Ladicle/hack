package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/golang/glog"
	"github.com/kyokomi/emoji"
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
		fmt.Sprintf("config file (default ~/%v)", config.DefaultConfig))

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
	err := config.Load(cfgFile)
	if err == nil {
		return
	}
	if !os.IsNotExist(err) {
		glog.Fatal(err)
	}

	// set defaults at first time
	emoji.Println(":tada:Welcome to the Hack!\n\nBefore start hacking, please answer some questions.\n")

	count := 1
	if err := initBaseDir(count); err != nil {
		glog.Fatal(err)
	}
	count++
}

func initBaseDir(count int) error {
	baseDir := config.BaseDir()
	if baseDir != "" {
		return nil
	}

	emoji.Printf("%v. Where do you put contests code? (default: ~/%v)\n->",
		count, config.DefaultBaseDir)
	fmt.Scanf("%s", &baseDir)
	baseDir = strings.TrimSpace(baseDir)

	u, err := user.Current()
	if err != nil {
		return err
	}
	if baseDir == "" {
		baseDir = filepath.Join(u.HomeDir, config.DefaultBaseDir)
	}
	if strings.HasPrefix(baseDir, "~") {
		baseDir = strings.Replace(baseDir, "~", u.HomeDir, 1)
	}
	config.SetBaseDir(baseDir)
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
