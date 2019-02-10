package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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

var (
	startNumber int
	interactive bool
)

func (c *sampleCmd) parse(args []string) error {
	fmt.Println(args)
	os.Args = args
	flag.BoolVar(&interactive, "i", false, "Tell you an answer if you continue input samples. (default: false)")
	flag.Parse()

	if config.C.CurrentQuizz == "" {
		return fmt.Errorf("could not get a current quiz: please set a current quiz (ex. hack jump)")
	}

	if flag.NArg() > 0 {
		i, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			return fmt.Errorf("%q is invalid sample number", flag.Arg(0))
		}
		startNumber = i
	} else {
		startNumber = nextSampleNum()
	}
	return nil
}

func (c *sampleCmd) run(args []string, opt Option) error {
	if err := c.parse(args); err != nil {
		return err
	}

	for n := startNumber; ; n++ {
		inSample := fmt.Sprintf("%d.in", n)
		outSample := fmt.Sprintf("%d.out", n)

		if _, err := os.Stat(inSample); err == nil {
			msg := fmt.Sprintf("%q is already exists. Skip it?", inSample)
			if yes, err := ansIsY(msg, c.IO); err != nil {
				return err
			} else if yes {
				if _, err := fmt.Fprintln(c.IO); err != nil {
					return err
				}
				continue
			}
		} else if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("unexpected error occurred: %v", err)
		}

		if _, err := fmt.Fprintf(c.IO, "%v:\n", inSample); err != nil {
			return err
		}
		if err := readAndCreateFile(genPathInQuizDir(inSample)); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(c.IO, "%v:\n", outSample); err != nil {
			return err
		}
		if err := readAndCreateFile(genPathInQuizDir(outSample)); err != nil {
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

func nextSampleNum() int {
	fs, err := ioutil.ReadDir(genQuizDir())
	if err != nil {
		return 1
	}
	var flag bool
	for _, f := range fs {
		n := strings.Split(f.Name(), ".")
		if len(n) < 2 {
			continue
		}
		if n[1] == "in" {
			if flag {
				i, _ := strconv.Atoi(n[0])
				return i
			}
			flag = true
		}
	}
	return 1
}

func genQuizDir() string {
	return filepath.Join(
		config.C.Contest.Path,
		config.C.CurrentQuizz)
}

func genPathInQuizDir(name string) string {
	return filepath.Join(genQuizDir(), name)
}

func ansIsY(msg string, io io.Writer) (bool, error) {
	var ans string
	if _, err := fmt.Fprintf(io, "%s (y/n): ", msg); err != nil {
		return false, err
	}
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
		return fmt.Errorf("failed to open file %q: %v", path, err)
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
