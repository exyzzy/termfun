package termfun

// canvas.go supports simple drawing and plotting using unicode block characters
// StringDense renders each character space as a 2x2 set of pixels
// StringAspect renders each character space as a 1x2 set of pixels which are more like square

import (
	"fmt"
	"strings"
)

// Canvas holds a bit buffer for plotting
// Canvas assumes a virtual canvas of bits arranged as
// width 0 ->
// height
// 0
// |
// v
// Upper left is 0,0
// bits are arranged within a byte as:
//
//	x x
// y 1 2
// y 4 8
// (y, x)
// (00)(01)
// (10)(11)
// if xor is true the new plot is xor'd on the existing canvas, otherwise it is just set
// if wrap is true the x, y coordinates are moded to wrap to the boundary
type Canvas struct {
	pwidth, pheight int //pixel
	bwidth, bheight int //byte
	value           *[]byte
	xor             bool
	wrap            bool
}

// Return a new initialized Canvas
func NewCanvas(w, h int) *Canvas {
	c := Canvas{}
	c.Init(w, h)
	return &c
}

// Init an existing Canvas, w and h are in pixels
// So Init(80,80) renders an array of 40x40 characters as StringDense
// and an array of 80x40 characters as StringAspect
func (c *Canvas) Init(w, h int) {
	c.xor = false
	c.wrap = false
	c.pwidth = w
	c.pheight = h
	c.bwidth = w/2 + 1
	c.bheight = h/2 + 1
	myarr := make([]byte, c.bwidth*c.bheight)
	c.value = &myarr
}

// Helpers

// Width returns the width of the canvas, in pixels
func (c *Canvas) Width() int { return c.pwidth }

// Height returns the height of the canvas, in pixels
func (c *Canvas) Height() int { return c.pheight }

// PlotXor sets the plot mode to xor the new pixel on the existing pixel
func (c *Canvas) PlotXor() { c.xor = true }

// PlotOr sets the plot mode to or the new pixel on the existing pixel (default)
func (c *Canvas) PlotOr() { c.xor = false }

// Plotwrap sets the plot mode to wrap cordinates on the borders
func (c *Canvas) PlotWrap() { c.wrap = true }

// Plotwrap sets the plot mode to not wrap cordinates on the borders (default)
func (c *Canvas) PlotUnwrap() { c.wrap = false }

// Plot a pixel on the Canvas
func (c *Canvas) Plot(x, y int) {
	if c.wrap {
		x = (x + c.pwidth) % c.pwidth
		y = (y + c.pheight) % c.pheight
	}
	if x < 0 {
		x = 0
	}
	if x >= c.pwidth {
		x = c.pwidth - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= c.pheight {
		y = c.pheight - 1
	}
	xi := x / 2
	xm := x % 2
	yi := y / 2
	ym := y % 2
	index := yi*c.bwidth + xi
	bit := byte(0x1 << (ym*2 + xm))
	// fmt.Println(x, y, " = ", xi, xm, "-", yi, ym, " > ", index, bit)
	if c.xor {
		(*c.value)[index] ^= bit
	} else {
		(*c.value)[index] |= bit
	}
}

// Read a location from the Canvas to see if it is set
func (c *Canvas) Read(x, y int) bool {
	if c.wrap {
		x = (x + c.pwidth) % c.pwidth
		y = (y + c.pheight) % c.pheight
	}
	if x < 0 {
		x = 0
	}
	if x >= c.pwidth {
		x = c.pwidth - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= c.pheight {
		y = c.pheight - 1
	}
	xi := x / 2
	xm := x % 2
	yi := y / 2
	ym := y % 2
	index := yi*c.bwidth + xi
	bit := byte(0x1 << (ym*2 + xm))
	return (*c.value)[index]&bit == bit
}

// Clear the Canvas
func (c *Canvas) Clear() {
	for i := range *c.value {
		(*c.value)[i] = byte(0)
	}
}

// Render the Canvas to a string as "square" (1x2) pixels
func (c *Canvas) StringAspect() string {
	var str string
	for y := 0; y < c.bheight; y++ {
		for x := 0; x < c.bwidth; x++ {
			index := y*c.bwidth + x
			str += fmt.Sprintf("%c%c", BlocksAspect[(*c.value)[index]][0], BlocksAspect[(*c.value)[index]][1])
		}
		str += "\r\n"
	}
	return str
}

// Render the Canvas to a string as "square" (1x2) pixels, with a border
func (c *Canvas) StringAspectBorder() string {
	var str string
	str = fmt.Sprintf("%c", SBox_UL)
	str += strings.Repeat(fmt.Sprintf("%c", SBox_Horiz), c.bwidth*2)
	str += fmt.Sprintf("%c\r\n", SBox_UR)
	for y := 0; y < c.bheight; y++ {
		str += fmt.Sprintf("%c", SBox_Vert)
		for x := 0; x < c.bwidth; x++ {
			index := y*c.bwidth + x
			str += fmt.Sprintf("%c%c", BlocksAspect[(*c.value)[index]][0], BlocksAspect[(*c.value)[index]][1])
		}
		str += fmt.Sprintf("%c\r\n", SBox_Vert)
	}
	str += fmt.Sprintf("%c", SBox_LL)
	str += strings.Repeat(fmt.Sprintf("%c", SBox_Horiz), c.bwidth*2)
	str += fmt.Sprintf("%c\r\n", SBox_LR)
	return str
}

// Render the Canvas to a string as dense (2x2) rectangle pixels
func (c *Canvas) StringDense() string {
	var str string
	for y := 0; y < c.bheight; y++ {
		for x := 0; x < c.bwidth; x++ {
			index := y*c.bwidth + x
			str += fmt.Sprintf("%c", BlocksDense[(*c.value)[index]])
		}
		str += "\r\n"
	}
	return str
}

// Render the Canvas to a string as dense (2x2) rectangle pixels, with a border
func (c *Canvas) StringDenseBorder() string {
	var str string
	str = fmt.Sprintf("%c", SBox_UL)
	str += strings.Repeat(fmt.Sprintf("%c", SBox_Horiz), c.bwidth)
	str += fmt.Sprintf("%c\r\n", SBox_UR)
	for y := 0; y < c.bheight; y++ {
		str += fmt.Sprintf("%c", SBox_Vert)
		for x := 0; x < c.bwidth; x++ {
			index := y*c.bwidth + x
			str += fmt.Sprintf("%c", BlocksDense[(*c.value)[index]])
		}
		str += fmt.Sprintf("%c\r\n", SBox_Vert)
	}
	str += fmt.Sprintf("%c", SBox_LL)
	str += strings.Repeat(fmt.Sprintf("%c", SBox_Horiz), c.bwidth)
	str += fmt.Sprintf("%c\r\n", SBox_LR)
	return str
}

// 1 2
// 4 8
var BlocksDense = [16]int{
	0x20,   // 0 all empty
	0x2598, // 1 UL
	0x259d, // 2 UR
	0x2580, // 3 UL, UR
	0x2596, // 4 LL
	0x258c, // 5 UL, LL
	0x259e, // 6 UR, LL
	0x259b, // 7 UL, UR, LL
	0x2597, // 8 LR
	0x259a, // 9 UL, LR
	0x2590, // 10 UR, LR
	0x259c, // 11 UL, UR, LR
	0x2584, // 12 LL, LR
	0x2599, // 13 UL, LL, LR
	0x259f, // 14 UR, LL, LR
	0x2588, // 15 all
}

var BlocksAspect = [16][2]int{
	{0x20, 0x20},     // 0 all empty
	{0x2580, 0x20},   // 1 UL
	{0x20, 0x2580},   // 2 UR
	{0x2580, 0x2580}, // 3 UL, UR
	{0x2584, 0x20},   // 4 LL
	{0x2588, 0x20},   // 5 UL, LL
	{0x2584, 0x2580}, // 6 UR, LL
	{0x2588, 0x2580}, // 7 UL, UR, LL
	{0x20, 0x2584},   // 8 LR
	{0x2580, 0x2584}, // 9 UL, LR
	{0x20, 0x2588},   // 10 UR, LR
	{0x2580, 0x2588}, // 11 UL, UR, LR
	{0x2584, 0x2584}, // 12 LL, LR
	{0x2588, 0x2584}, // 13 UL, LL, LR
	{0x2584, 0x2588}, // 14 UR, LL, LR
	{0x2588, 0x2588}, // 15 all
}

// Line uses Bresenham's algorithm to plot a line
func (c *Canvas) Line(x0, y0, x1, y1 int) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		c.Plot(x0, y0)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// Bmp plots a [][]byte to the Canvas
func (c *Canvas) Bmp(x, y int, bmp [][]byte) {
	for j, a := range bmp {
		for i, e := range a {
			for b := 0; b < 8; b++ {
				bit := byte(0x80 >> b)
				if e&bit == bit {
					c.Plot(x+(i*8)+b, y+j)
				}
			}
		}
	}
}
