package main

import (
	"fmt"
	"io"

	"github.com/mattn/go-runewidth"
)

func fprintLFW(o io.Writer, v interface{}, width int) {
	fmt.Fprint(o, runewidth.FillRight(fmt.Sprintf("%v", v), width))
}

func fprintRFW(o io.Writer, v interface{}, width int) {
	fmt.Fprint(o, runewidth.FillLeft(fmt.Sprintf("%v", v), width))
}
