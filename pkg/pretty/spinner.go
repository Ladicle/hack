package pretty

import (
	"fmt"
	"io"
	"time"

	"github.com/logrusorgru/aurora"
)

type Progress struct {
	spinner *Spinner
	writer  io.Writer
}

func NewProgress(w io.Writer) *Progress {
	return &Progress{
		spinner: NewSpinner(w),
		writer:  w,
	}
}

func (p *Progress) Start(msg string) {
	p.spinner.Stop(true)
	p.spinner.Start(msg)
}

func (p *Progress) End(success bool) {
	p.spinner.Stop(success)
}

// custom CLI loading spinner for kind
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
