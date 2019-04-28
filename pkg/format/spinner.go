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
	stop   chan interface{}
	ticker *time.Ticker
	writer io.Writer
}

func NewSpinner(w io.Writer) *Spinner {
	return &Spinner{
		frames: spinnerFrames,
		stop:   make(chan interface{}, 1),
		ticker: time.NewTicker(time.Millisecond * 100),
		writer: w,
	}
}

func (s *Spinner) Start(msg string) {
	go func() {
		for {
			for _, frame := range s.frames {
				select {
				case <-s.stop:
					fmt.Fprint(s.writer, "\r")
					fmt.Fprintf(s.writer, " %v %s\n",
						aurora.Green("✓").Bold(), msg)
					// if success {
					// 	fmt.Fprintf(s.writer, " %v %s\n",
					// 		aurora.Green("✓").Bold(), msg)
					// } else {
					// 	fmt.Fprintf(s.writer, " %v %s\n",
					// 		aurora.Red("✗").Bold(), msg)
					// }
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
