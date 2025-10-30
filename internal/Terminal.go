package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

func Init() *term.State {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		//?log
		return nil
	}
	return oldState
}

func Deinit(oldState *term.State) {
	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
}

func ClearLine() {
	fmt.Printf("\033[2K")
}

func ClearLineEnd() {
	fmt.Printf("\033[0K")
}

func ClearLineStart() {
	fmt.Printf("\033[1K")
}

func CursorHide() {
	fmt.Printf("\033[?25l")
}

func CursorPositionRestore() {
	fmt.Printf("\033[u")
}

func CursorSavePosition() {
	fmt.Printf("\033[s")
}

func CursorShow() {
	fmt.Printf("\033[?25h")
}

func MoveTo(x, y int) {
	fmt.Printf("\033[%d;%df", y, x)
}

// MoveDown moves the cursor to the beginning of n lines down.
func MoveDown(n int) {
	fmt.Printf("\033[%dE", n)
}

// MoveUp moves the cursor to the beginning of n lines up.
func MoveUp(n int) {
	fmt.Printf("\033[%dF", n)
}

func ScreenClear() {
	fmt.Printf("\033[2J")
	MoveTo(1, 1)
}

func ScreenRestore() {
	fmt.Printf("\033[?47l")
}

func ScreenSave() {
	fmt.Printf("\033[?47h")
}

func ScreenSize() (w int, h int) {
	fmt.Printf("\033[18t") // Report Size
	// The response will be in the form ESC [ 8 ; height ; width t
	// read from stdin to get the response
	buf := ""
	err := error(nil)
	for {
		var charBuf [1]byte
		_, err = os.Stdin.Read(charBuf[:])
		if err != nil {
			break
		}
		b := charBuf[0]
		if b == 't' {
			break
		}

		buf += string(b)
	}
	if err != nil || len(buf) < 6 || buf[0] != '\033' || buf[1] != '[' || buf[2] != '8' || buf[3] != ';' {
		// ?log?
		return 80, 25 // default size
	}
	// parse height and width
	hw := strings.Split(buf[4:], ";")
	if len(hw) != 2 {
		return 80, 25
	}
	w, err = strconv.Atoi(hw[1])
	if err != nil {
		return 80, 25
	}

	h, err = strconv.Atoi(hw[0])
	if err != nil {
		return 80, 25
	}

	return w, h
}

/*
ANSI codes for screen management
Save the current screen: ESC[?47h
This saves the content of the current screen buffer to an alternate buffer.
Restore the saved screen: ESC[?47l
This restores the content from the alternate buffer, effectively "capturing" and then "recapturing" the screen.
Clear the entire screen: ESC[2J
This code clears the screen content but does not save it to an alternate buffer.

smcup
\E7 saves the cursor's position
\E[?47h switches to the alternate screen
rmcup
\E[2J clears the screen (assumed to be the alternate screen)
\E[?47l switches back to the normal screen
\E8 restores the cursor's position.

The syntax is ESC[?1049h and ESC[?1049l -- here h activates the alternate buffer and l returns to normal mode. Other similar codes are 1047 and 1048. But apparently they are older and have less functionality.


https://en.wikipedia.org/wiki/ANSI_escape_code
https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
https://xtermjs.org/docs/api/vtfeatures/
https://invisible-island.net/xterm/ctlseqs/ctlseqs.html
https://github.com/leaanthony/go-ansi-parser

Xterm allows the window title to be set by ESC ]0;this is the window title BEL

*/
