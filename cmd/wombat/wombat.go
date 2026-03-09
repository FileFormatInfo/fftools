package main

import (
	"fmt"
	"os"
	"time"

	. "github.com/FileFormatInfo/fftools/internal"
	"github.com/spf13/pflag"
)

var (
	BUILDER = "unknown"
	COMMIT  = "(local)"
	LASTMOD = "(local)"
	VERSION = "internal"
)

func main() {

	var version = pflag.Bool("version", false, "Version info")

	pflag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "wombat version %s (built by %s on %s, commit %s)\n", VERSION, BUILDER, LASTMOD, COMMIT)
		return
	}

	oldState := Init()
	defer Deinit(oldState)

	ScreenSave()
	w, h := ScreenSize()

	xMid := float32(w) / 2
	yMid := float32(h) / 2

	steps := float32(17)

	xStep := float32(xMid) / float32(steps)
	yStep := float32(yMid) / float32(steps)

	CursorSavePosition()
	CursorHide()
	ScreenClear()
	for loop := float32(0); loop < steps; loop++ {
		MoveTo(int(xMid+loop*xStep), int(yMid+loop*yStep))
		fmt.Printf("*")
		MoveTo(int(xMid+loop*xStep), int(yMid-loop*yStep))
		fmt.Printf("*")
		MoveTo(int(xMid-loop*xStep), int(yMid-loop*yStep))
		fmt.Printf("*")
		MoveTo(int(xMid-loop*xStep), int(yMid+loop*yStep))
		fmt.Printf("*")

		time.Sleep(20 * time.Millisecond)
		ScreenClear()
	}

	CursorPositionRestore()
	CursorShow()
	ScreenRestore()
}
