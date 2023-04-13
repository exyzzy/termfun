// Example of unicode block characters to draw graphics
package main

import (
	"fmt"
	"math"
	"os"

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

	c := termfun.NewCanvas(80, 50)
	CanvasTest(c)
	fmt.Print(c.StringDenseBorder(), "\r\n")
	fmt.Print(c.StringAspectBorder(), "\r\n")

}

// test by plotting sine fn and bmp chars
func CanvasTest(c *termfun.Canvas) {
	var i0, j0, i1, j1 int
	c.Line(c.Width()/2, 0, c.Width()/2, c.Height()) //vertical axis
	for i := 0; i < c.Width(); i++ {
		c.Plot(i, c.Height()/2) // horizontal axis
		//plot sine
		y := math.Sin(float64(i) / float64(c.Width()) * math.Pi * 2)
		j := int((y + 1.0) / 2 * float64(c.Height()))
		// c.Plot(i, j) //could, but below with c.Line is better
		if i == 0 {
			i0 = i
			j0 = j
		} else {
			i1 = i
			j1 = j
			c.Line(i0, j0, i1, j1)
			i0 = i1
			j0 = j1
		}
	}

	// Use c.Bmp to render "Go" as pixels
	var b = [][]byte{
		{0x00, 0x00},
		{0x00, 0x00},
		{0x60, 0x00},
		{0x90, 0x60},
		{0x80, 0x90},
		{0x90, 0x90},
		{0x70, 0x60},
		{0x00, 0x00}}
	c.Bmp(c.Width()/5, 0, b)
}
