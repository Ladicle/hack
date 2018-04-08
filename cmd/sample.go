package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
)

// NewSampleCmd samples contest samplermation.
func NewSampleCmd(io io.Writer) Command {
	s := sampleCmd{IO: io}
	return Command{
		Name:        "sample",
		Short:       "sample [-i] [NUMBER]",
		Description: "Sample creates input/output sample",
		Run:         s.run,
	}
}

type sampleCmd struct {
	IO io.Writer
}

func (c *sampleCmd) run(args []string, opt Option) error {
	var start int
	var interactive bool

	flag.BoolVar(&interactive, "i", false, "Tell you an answer if you continue input samples. (default: false)")
	flag.Parse()

	switch flag.NArg() {
	case 0:
		start = 1
	case 1:
		if i, err := strconv.Atoi(flag.Arg(0)); err != nil {
			fmt.Errorf("%q is not number", flag.Arg(0))
		} else {
			start = i
		}
	default:
		return fmt.Errorf("invalid number of arguments")
	}

	for n := start; ; n++ {
		inSample := fmt.Sprintf("%d.in", n)
		outSample := fmt.Sprintf("%d.out", n)

		if _, err := os.Stat(inSample); err == nil {
			msg := fmt.Sprintf("%q is already exists. Skip it?", inSample)
			if yes, err := ansIsY(msg, c.IO); err != nil {
				return err
			} else if yes {
				fmt.Fprintln(c.IO)
				continue
			}
		} else if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("unexpected error occurred: %v", err)
		}

		if config.C.CurrentQuizz == "" {
			return fmt.Errorf("could not get a current quiz: please set a current quiz (ex. hack jump)")
		}

		fmt.Fprintf(c.IO, "%v:\n", inSample)
		if err := readAndCreateFile(genQuizDir(inSample)); err != nil {
			return err
		}

		fmt.Fprintf(c.IO, "%v:\n", outSample)
		if err := readAndCreateFile(genQuizDir(outSample)); err != nil {
			return err
		}

		if !interactive {
			continue
		}

		if yes, err := ansIsY("Continue to create sample?", c.IO); err != nil {
			return err
		} else if !yes {
			break
		}
	}
	return nil
}

func genQuizDir(name string) string {
	return filepath.Join(
		config.C.Contest.Path,
		config.C.CurrentQuizz,
		name)
}

func ansIsY(msg string, io io.Writer) (bool, error) {
	var ans string
	fmt.Fprintf(io, "%s (y/n): ", msg)
	if _, err := fmt.Scanf("%s", &ans); err != nil {
		return false, fmt.Errorf("failed to read answer: %v", err)
	}
	if ans == "y" {
		return true, nil
	}
	return false, nil
}

func readAndCreateFile(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Errorf("failed to open file %q: %v", path, err)
	}
	defer f.Close()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(f)
	defer out.Flush()

	for {
		s, err := readNextLine(in)
		if err != nil {
			return err
		}
		if strings.TrimSpace(s) == "" {
			break
		}
		if _, err = out.WriteString(fmt.Sprintf("%s\n", s)); err != nil {
			return err
		}
	}
	return nil
}

func readNextLine(in *bufio.Reader) (string, error) {
	var b []byte
	for {
		l, pre, err := in.ReadLine()
		if err != nil {
			return "", err
		}
		b = append(b, l...)
		if !pre {
			break
		}
	}
	return string(b), nil
}
