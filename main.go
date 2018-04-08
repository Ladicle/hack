package main

import (
	"flag"
	"fmt"
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

func main() {
	io, ioErr := os.Stdout, os.Stderr
	cmd.LoadCmd(io)

	flag.StringVar(&cmd.ConfigPath, "config", "", "")
	flag.StringVar(&cmd.ConfigPath, "c", "", "")

	flag.StringVar(&cmd.OutputDirectory, "output", "", "")
	flag.StringVar(&cmd.OutputDirectory, "o", "", "")

	flag.Usage = func() {
		fmt.Fprintf(io, `Usage: hack [OPTIONS] COMMAND

Options:
  -c --config         Configuration path (default: ~/%s)
  -o --output         Output directory (default: ~/%s)
  -h --help           Show this help message

Commands:\n`, defaultConfigPath, defaultOutputDirectory)
		cmd.PrintUsage(io)
	}

	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(ioErr, "Invalid number of arguments")
		os.Exit(1)
	}

	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintln(ioErr, "Could not find HOME")
		os.Exit(1)
	}
	if cmd.ConfigPath == "" {
		cmd.ConfigPath = filepath.Join(home, defaultConfigPath)
	}
	if cmd.OutputDirectory == "" {
		cmd.OutputDirectory = filepath.Join(home, defaultOutputDirectory)
	}

	if err := config.LoadConfig(cmd.ConfigPath); err != nil {
		fmt.Fprintf(ioErr, "Filed to load configuration from %v\n", cmd.ConfigPath)
	}

	os.Args = flag.Args()
	if err := cmd.HandleCmd(flag.Arg(0), flag.Args(), cmd.Option{}); err != nil {
		fmt.Fprintln(ioErr, err)
		os.Exit(1)
	}
}
