package format

import (
	"fmt"
	"io"

	"github.com/kyokomi/emoji/v2"
	"github.com/logrusorgru/aurora/v3"
)

type Progress struct {
	started bool
	msg     string
	writer  io.Writer
}

func NewProgress(w io.Writer) *Progress {
	return &Progress{writer: w}
}

func (p *Progress) StartWithEmojiMsg(emojiMsg string) {
	p.Start(emoji.Sprintf(emojiMsg))
}

func (p *Progress) Start(msg string) {
	p.End(true)

	fmt.Fprintf(p.writer, " - %s", msg)

	p.msg = msg
	p.started = true
}

func (p *Progress) End(success bool) {
	if !p.started {
		return
	}

	var icon aurora.Value
	if success {
		icon = aurora.Green("✓").Bold()

	} else {
		icon = aurora.Red("✗").Bold()
	}
	fmt.Fprintf(p.writer, "\r %v %s\n", icon, p.msg)
}
