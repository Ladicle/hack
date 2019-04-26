package format

import (
	"io"
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
