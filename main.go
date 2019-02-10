package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Ladicle/hack/cmd"
	"github.com/Ladicle/hack/pkg/config"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	defaultConfigPath      = ".hack"
	defaultOutputDirectory = "contest"
)

func init() {
	cmd.LoadCmd(os.Stdout)

	flag.StringVar(&cmd.ConfigPath, "config", "", "")
	flag.StringVar(&cmd.ConfigPath, "c", "", "")
	flag.StringVar(&cmd.OutputDirectory, "output", "", "")
	flag.StringVar(&cmd.OutputDirectory, "o", "", "")
	flag.Usage = func() {
		fmt.Printf(`Usage: hack [OPTIONS] COMMAND

Options:
  -c --config         Configuration path (default: ~/%s)
  -o --output         Output directory (default: ~/%s)
  -h --help           Show this help message

Commands:
`, defaultConfigPath, defaultOutputDirectory)
		cmd.PrintUsage()
	}
}

func main() {
	flag.Parse()
	if err := validation(); err != nil {
		log.Fatalf("Failed a root command validation: %v", err)
	}

	if err := config.LoadConfig(cmd.ConfigPath); err != nil {
		log.Fatalf("Filed to load configuration from %v\n", cmd.ConfigPath)
	}

	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func validation() error {
	if flag.NArg() == 0 {
		return fmt.Errorf("no arguments are specified")
	}

	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("not found home directory")
	}
	if cmd.ConfigPath == "" {
		cmd.ConfigPath = filepath.Join(home, defaultConfigPath)
	}
	if cmd.OutputDirectory == "" {
		cmd.OutputDirectory = filepath.Join(home, defaultOutputDirectory)
	}
	return nil
}

func run() error {
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("filed to get working directory because %v", err)
	}
	if err := cmd.HandleCmd(flag.Arg(0), flag.Args()[1:], cmd.Option{WorkDir: workDir}); err != nil {
		return err
	}
	return nil
}
