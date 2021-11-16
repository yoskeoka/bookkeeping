package main

import (
	"fmt"
	"io"

	"github.com/mattn/go-runewidth"
)

func fprintLFW(o io.Writer, s string, width int) {
	fmt.Fprint(o, runewidth.FillRight(s, width))
}
