package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// These variables are set at the build time.
	version string
	gitRepo string
)

func main() {
	io, ioErr := os.Stdout, os.Stderr

	flag.Usage = func() {
		fmt.Fprintln(io, `Usage: hack [OPTIONS] COMMAND

Options:
  -h --help      Show this help message

Commands:
  atcoder        Support AtCoder contests
  version        Show this command version`)
	}

	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(ioErr, "Invalid number of arguments")
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "atcoder":
		fmt.Fprintln(ioErr, "Not implemented")
		os.Exit(1)
	case "version":
		fmt.Fprintf(io, "%v version is %v", gitRepo, version)
		os.Exit(0)
	default:
		fmt.Fprintln(ioErr, "Unknown command")
		os.Exit(1)
	}
}
