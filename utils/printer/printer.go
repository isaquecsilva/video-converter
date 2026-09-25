package printer

import (
	"cmp"
	"fmt"
	"io"
	"os"
)

type Printer interface {
	Info(ln bool, msg string, args ...any)
	Success(ln bool, msg string, args ...any)
	Warn(ln bool, msg string, args ...any)
	Error(ln bool, msg string, args ...any)
}

type defaultPrinter struct {
	w         io.Writer
	colorsMap ColorMap
}

func NewDefaultPrinter(w io.Writer, colorsMap ColorMap) Printer {
	if colorsMap == nil {
		colorsMap = defaultColorMap()
	}

	FillMissingColorMapColors(colorsMap)

	return defaultPrinter{
		w:         cmp.Or[io.Writer](w, os.Stdout),
		colorsMap: colorsMap,
	}
}

func (d defaultPrinter) Info(ln bool, msg string, args ...any) {
	d.colorsMap["info"].Fprintf(d.w, msg, args...)
	d.println(ln)
}

func (d defaultPrinter) Infoln(ln bool, msg string, args ...any) {
	d.colorsMap["info"].Fprintf(d.w, msg, args...)
	d.println(ln)
}

func (d defaultPrinter) Success(ln bool, msg string, args ...any) {
	d.colorsMap["success"].Fprintf(d.w, msg, args...)
	d.println(ln)
}

func (d defaultPrinter) Warn(ln bool, msg string, args ...any) {
	d.colorsMap["warn"].Fprintf(d.w, msg, args...)
	d.println(ln)
}

func (d defaultPrinter) Error(ln bool, msg string, args ...any) {
	d.colorsMap["error"].Fprintf(d.w, msg, args...)
	d.println(ln)
}

func (d defaultPrinter) println(ln bool) {
	if ln {
		fmt.Println()
	}
}
