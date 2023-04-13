// Conway's Game of Life.
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/exyzzy/termfun"
)

func main() {
	c := termfun.NewCanvas(80, 50)
	c.PlotWrap()
	Randomize(c)
	fmt.Print(termfun.ED(termfun.EraseAll), "\r\n")

	for i := 0; i < 400; i++ {
		fmt.Print(termfun.CUP(0, 0))
		fmt.Print(c.StringDenseBorder(), "\r\n")
		fmt.Printf("Iteration: %d\r\n", i)
		c = NextFrame(c)
		time.Sleep(time.Second / 10)
	}
}

// set some random pixels
func Randomize(c *termfun.Canvas) {
	for i := 0; i < (c.Width() * c.Height() / 5); i++ {
		c.Plot(rand.Intn(c.Width()), rand.Intn(c.Height()))
	}
}

// calculate the next frame
func NextFrame(c *termfun.Canvas) *termfun.Canvas {
	cnext := termfun.NewCanvas(c.Width(), c.Height())
	cnext.PlotWrap()
	for y := 0; y < c.Height(); y++ {
		for x := 0; x < c.Width(); x++ {
			if PixelInNext(x, y, c) {
				cnext.Plot(x, y)
			}

		}
	}
	return cnext
}

// check if a pixel is set in the next frame
func PixelInNext(x, y int, c *termfun.Canvas) bool {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && c.Read(x+i, y+j) {
				count++
			}
		}
	}
	return count == 3 || count == 2 && c.Read(x, y)
}
