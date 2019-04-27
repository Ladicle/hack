package format

import (
	"io"
	"os"
)

type HackRobot struct {
	Progress *Progress
	Robot    *Robot
}

func NewHackRobot(w io.Writer) *HackRobot {
	return &HackRobot{
		Progress: NewProgress(w),
		Robot:    NewRobot(w),
	}
}

// Fatal prints error message to the Stderr and calls Exit(1).
func (h *HackRobot) Fatal(format string, args ...interface{}) {
	h.Progress.End(false)
	h.Robot.FprintlnRed(os.Stderr, format, args...)
	os.Exit(1)
}

// Info prints information.
func (h *HackRobot) Info(format string, args ...interface{}) {
	h.Robot.PrintlnGreen(format, args...)
}

// Start starts progress spinner and shows message with emoji.
func (h *HackRobot) Start(msg string) {
	h.Progress.StartWithEmojiMsg(msg)
}

// End finalize robot and progress.
func (h *HackRobot) End() {
	h.Progress.End(true)
}
