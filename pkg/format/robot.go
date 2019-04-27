package format

import (
	"fmt"
	"io"

	"github.com/kyokomi/emoji"
	"github.com/logrusorgru/aurora"
)

type Robot struct {
	icon   string
	writer io.Writer
	count  int
}

func NewRobot(w io.Writer) *Robot {
	return &Robot{
		icon:   ":robot:",
		writer: w,
	}
}

// Printls prints lines.
func (r *Robot) Printls(msg ...string) {
	var ls string
	for _, m := range msg {
		ls += r.getIcon() + "< " + m + "\n"
		r.count++
	}
	fmt.Fprint(r.writer, ls)
}

// PrintlnYellow prints formatted line string with read arguments.
func (r *Robot) PrintlnYellow(format string, args ...interface{}) {
	var aargs []interface{}
	for _, a := range args {
		aargs = append(aargs, aurora.Yellow(a).Bold())
	}
	r.printlnf(format, aargs...)
}

// PrintlnGreen prints formatted line string with read arguments.
func (r *Robot) PrintlnGreen(format string, args ...interface{}) {
	var aargs []interface{}
	for _, a := range args {
		aargs = append(aargs, aurora.Green(a).Bold())
	}
	r.printlnf(format, aargs...)
}

// PrintlnRed prints formatted line string with read arguments.
func (r *Robot) PrintlnRed(format string, args ...interface{}) {
	var aargs []interface{}
	for _, a := range args {
		aargs = append(aargs, aurora.Red(a).Bold())
	}
	r.printlnf(format, aargs...)
}

func (r *Robot) printlnf(format string, args ...interface{}) {
	emoji.Fprintf(r.writer, r.getIcon()+"< "+format+"\n", args...)
	r.count++
}

func (r *Robot) getIcon() string {
	if r.count == 0 {
		return r.icon
	}
	return "   "
}
