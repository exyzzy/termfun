// Example of decoding CSI codes from keyboard special keys
package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
	"unicode/utf16"

	"github.com/exyzzy/termfun"
	"golang.org/x/term"
)

func main() {
	in := os.Stdin
	oldState, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(in.Fd()), oldState)

	reader := bufio.NewReader(in)

	var k rune
	fmt.Print("Press keys (Ctrl-C to exit)..\r\n")
	for {
		// k, _, err = reader.ReadRune() //normally do this but if arrow keys are needed then use ReadKey
		k, _, err = termfun.ReadKey(reader)
		if err != nil {
			panic(err)
		}
		if unicode.IsControl(k) || utf16.IsSurrogate(k) {
			if txt, ok := keyMap[k]; ok {
				fmt.Printf("%d (%s)\r\n", k, txt)
			} else {
				fmt.Printf("%d\r\n", k)
			}
		} else {
			fmt.Printf("%d ('%c')\r\n", k, k)
		}
		if k == termfun.CtrlC {
			break
		}
	}
}

var keyMap = map[rune]string{
	termfun.KeyUnknown: "Unknown",
	termfun.KeyUp:      "Up",
	termfun.KeyDown:    "Down",
	termfun.KeyLeft:    "Left",
	termfun.KeyRight:   "Right",
	termfun.KeyBackTab: "BackTab",
	termfun.KeyDel:     "Delete",
}
