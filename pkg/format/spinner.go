package format

import (
	"fmt"
	"io"
	"time"

	"github.com/logrusorgru/aurora"
)

var spinnerFrames = []string{
	"⠈⠁",
	"⠈⠑",
	"⠈⠱",
	"⠈⡱",
	"⢀⡱",
	"⢄⡱",
	"⢄⡱",
	"⢆⡱",
	"⢎⡱",
	"⢎⡰",
	"⢎⡠",
	"⢎⡀",
	"⢎⠁",
	"⠎⠁",
	"⠊⠁",
}

type Spinner struct {
	frames []string
	stop   chan bool
	ticker *time.Ticker
	writer io.Writer
}

func NewSpinner(w io.Writer) *Spinner {
	return &Spinner{
		frames: spinnerFrames,
		stop:   make(chan bool, 1),
		ticker: time.NewTicker(time.Millisecond * 100),
		writer: w,
	}
}

func (s *Spinner) Start(msg string) {
	go func() {
		for {
			for _, frame := range s.frames {
				select {
				case success := <-s.stop:
					fmt.Fprint(s.writer, "\r")
					if success {
						fmt.Fprintf(s.writer, " %v %s\n",
							aurora.Green("✓").Bold(), msg)
					} else {
						fmt.Fprintf(s.writer, " %v %s\n",
							aurora.Red("✗").Bold(), msg)
					}
					return
				case <-s.ticker.C:
					fmt.Fprintf(s.writer, "\r%s %s",
						aurora.Yellow(frame), msg)
				}
			}
		}
	}()
}

func (s *Spinner) Stop(success bool) {
	s.stop <- success
}
