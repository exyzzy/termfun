// Example of CSI codes to draw a box on the raw terminal
package main

import (
	"fmt"
	"os"
	"strings"

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

	str := termfun.ED(termfun.EraseAll) // clear screen
	str += Box(10, 10, 20, 20)
	fmt.Print(str, "\r\n")
	// fmt.Print(strconv.QuoteToASCII(str), "\r\n")
}

// draw a horizontal line - also in tile.go
func HLine(x1, x2, y, c int) string {
	return termfun.CUP(x1, y) + strings.Repeat(fmt.Sprintf("%c", c), x2-x1)
}

// draw a vertical line - also in tile.go
func VLine(x, y1, y2, c int) string {
	return termfun.CUP(x, y1) + strings.Repeat(fmt.Sprintf("%c%s%s", c, termfun.CUB(1), termfun.CUD(1)), y2-y1)
}

// place a character at a location - also in tile.go
func CharAt(x, y, c int) string {
	return termfun.CUP(x, y) + fmt.Sprintf("%c", c)
}

// draw box
func Box(x1, y1, x2, y2 int) string {
	var str string
	str = CharAt(x1, y1, termfun.SBox_UL)
	str += termfun.CUU(1)
	str += HLine(x1+1, x2, y1, termfun.SBox_Horiz)
	str += CharAt(x2, y1, termfun.SBox_UR)
	str += VLine(x2, y1+1, y2, termfun.SBox_Vert)
	str += CharAt(x2, y2, termfun.SBox_LR)
	str += HLine(x1+1, x2, y2, termfun.SBox_Horiz)
	str += CharAt(x1, y2, termfun.SBox_LL)
	str += VLine(x1, y1+1, y2, termfun.SBox_Vert)
	return str
}
