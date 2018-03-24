package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Ladicle/hack/cmd"
)

func main() {
	io, ioErr := os.Stdout, os.Stderr
	cmd.LoadCmd(io)

	flag.StringVar(&cmd.ConfigPath, "config", "", "")
	flag.StringVar(&cmd.ConfigPath, "c", "", "")

	flag.StringVar(&cmd.OutputDirectory, "output", "", "")
	flag.StringVar(&cmd.OutputDirectory, "o", "", "")

	flag.Usage = func() {
		fmt.Fprintln(io, `Usage: hack [OPTIONS] COMMAND

Options:
  -c --config         Configuration path (default: ~/.hack)
  -o --output         Output directory (default: ~/competition)
  -h --help           Show this help message

Commands:`)
		cmd.PrintUsage(io)
	}

	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(ioErr, "Invalid number of arguments")
		os.Exit(1)
	}

	os.Args = flag.Args()
	if err := cmd.HandleCmd(flag.Arg(0)); err != nil {
		fmt.Fprintln(ioErr, err)
		os.Exit(1)
	}
}
