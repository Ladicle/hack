package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Ladicle/hack/pkg/cmd/jump"
	"github.com/Ladicle/hack/pkg/cmd/list"
	"github.com/Ladicle/hack/pkg/config"
	"github.com/golang/glog"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
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
	rootCmd.AddCommand(list.NewCommand())
	rootCmd.AddCommand(jump.NewCommand())

	rootCmd.AddCommand(NewSampleCmd())
	rootCmd.AddCommand(NewCopyCmd())
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewSwitchCmd())
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
	emoji.Println("Welcome to the Hack!:tada:\n\n:robot:< I'm a hack bot!\n   < Before start hacking, please answer some questions.\n")

	count := 1
	if err := initBaseDir(count); err != nil {
		glog.Fatal(err)
	}
	count++

	if err := initAtCoder(count); err != nil {
		glog.Fatal(err)
	}
	count++
}

func initBaseDir(count int) error {
	var baseDir string

	fmt.Printf("%v. Where do you put contests code? (default: ~/%v)\n-> ",
		count, config.DefaultBaseDir)
	fmt.Scanln(&baseDir)
	fmt.Println()
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
	glog.V(4).Info("Saved base directory")
	return nil
}

func initAtCoder(count int) error {
	var ans string
	fmt.Printf("%v. Do you have AtCoder account? (y/n)\n-> ", count)
	fmt.Scanln(&ans)

	if ans != "y" {
		fmt.Printf("OK! I'll skip this.\n\n")
		return nil
	}

	var user, pass string
	fmt.Printf("%v.1. Tell me the username.\n-> ", count)
	fmt.Scanln(&user)
	fmt.Printf("%v.2. Tell me the password.\n-> ", count)
	bpass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	pass = string(bpass)
	fmt.Printf("\n\n")

	config.SetAtCoderUser(user)
	config.SetAtCoderPass(pass)

	glog.V(4).Info("Saved AtCoder username and password")
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
