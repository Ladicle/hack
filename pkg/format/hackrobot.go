package format

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

const (
	StateRuntimeError      = "RE"
	StateCompileError      = "CE"
	StateTimeLimitExceeded = "TLE"
	StateAnswerIsCorrect   = "AC"
	StateWrongAnswer       = "WA"
)

type HackRobot struct {
	Progress *Progress
	Robot    *Robot

	out    io.Writer
	errOut io.Writer
}

func NewHackRobot(w io.Writer) *HackRobot {
	return &HackRobot{
		Progress: NewProgress(w),
		Robot:    NewRobot(w),
		out:      w,
		errOut:   os.Stderr,
	}
}

func (h *HackRobot) PrettyDiff(s1, s2 string) {
	s1L := strings.Split(s1, "\n")
	s2L := strings.Split(s2, "\n")
	s1Len := len(s1L)

	for i, w := range s2L {
		if i >= s1Len {
			h.Printfln("%v\n%v", aurora.Red("<empty>"), aurora.Green(w))
			continue
		}
		if w == s1L[i] {
			h.Printfln(w)
		} else {
			h.Printfln("%v\n%v", aurora.Red(s1L[i]), aurora.Green(w))
		}
	}
	if s1Len > len(s2L) {
		for i := len(s2L); i < s1Len; i++ {
			h.Printfln("%v", aurora.Red(s1L[i]))
		}
	}
}

// Fatal prints error message to the Stderr and calls Exit(1).
func (h *HackRobot) Fatal(format string, args ...interface{}) {
	h.Robot.FprintlnRed(h.errOut, format, args...)
	os.Exit(1)
}

func (h *HackRobot) Printfln(format string, args ...interface{}) {
	fmt.Fprintf(h.out, format+"\n", args...)
}

// Info prints information.
func (h *HackRobot) Info(format string, args ...interface{}) {
	h.Robot.PrintlnGreen(format, args...)
}

func (h *HackRobot) State(state, msg string) {
	var v aurora.Value
	switch state {
	case StateAnswerIsCorrect:
		v = aurora.Green(state)
	case StateTimeLimitExceeded:
		v = aurora.Yellow(state)
	default:
		v = aurora.Red(state)
	}
	h.Robot.PrintlnYellow("[%v] %v", v, msg)
}

// Start starts progress spinner and shows message with emoji.
func (h *HackRobot) Start(format string, args ...interface{}) {
	h.Progress.StartWithEmojiMsg(fmt.Sprintf(format, args...))
}

// End finalize robot and progress.
func (h *HackRobot) Success() {
	h.Progress.End(true)
}

// Error stop progress and set error result.
func (h *HackRobot) Error() {
	h.Progress.End(false)
}
