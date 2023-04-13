package termfun

//csi.go supports the display control CSI codes

import (
	"fmt"
)

// ANSI CSI term codes
// https://en.wikipedia.org/wiki/ANSI_escape_code#CSI_(Control_Sequence_Introducer)_sequences

// Upper Left of screen is 1, 1

const CSI = "\x1b\x5b"

// CUU - Cursor Up
func CUU(n int) string {
	return (fmt.Sprintf("%s%dA", CSI, n))
}

// CUD - Cursor Down
func CUD(n int) string {
	return (fmt.Sprintf("%s%dB", CSI, n))
}

// CUF - Cursor Forward
func CUF(n int) string {
	return (fmt.Sprintf("%s%dC", CSI, n))
}

// CUB - Cursor Back
func CUB(n int) string {
	return (fmt.Sprintf("%s%dD", CSI, n))
}

// CNL - Cursor Next Line
func CNL(n int) string {
	return (fmt.Sprintf("%s%dE", CSI, n))
}

// CPL - Cursor Previous Line
func CPL(n int) string {
	return (fmt.Sprintf("%s%dF", CSI, n))
}

// CHA - Cursor Horizontal Absolute
func CHA(n int) string {
	return (fmt.Sprintf("%s%dG", CSI, n))
}

// CUP - Cursor Position
func CUP(m, n int) string {
	return (fmt.Sprintf("%s%d;%dH", CSI, n, m))
}

type EraseType int

const (
	EraseToEnd   EraseType = 0
	EraseToBegin EraseType = 1
	EraseAll     EraseType = 2
)

// ED - Erase in Display
// If n is 0, clear from cursor to end of screen. If n is 1, clear from cursor to beginning of the screen. If n is 2, clear entire screen (and moves cursor to upper left on DOS ANSI.SYS).
func ED(n EraseType) string {
	return (fmt.Sprintf("%s%dJ", CSI, n))
}

// EL - Erase in Line
// If n is 0 (or missing), clear from cursor to the end of the line. If n is 1, clear from cursor to beginning of the line. If n is 2, clear entire line. Cursor position does not change.
func EL(n EraseType) string {
	return (fmt.Sprintf("%s%dK", CSI, n))
}

// SU - Scroll Up
func SU(n int) string {
	return (fmt.Sprintf("%s%dS", CSI, n))
}

// SD - Scroll Down
func SD(n int) string {
	return (fmt.Sprintf("%s%dT", CSI, n))
}

// HVP - Horizontal Vertical Position
func HVP(m, n int) string {
	return (fmt.Sprintf("%s%d;%df", CSI, n, m))
}

type SGRType int

const (
	SGR_Off          SGRType = 0  // All attributes off
	SGR_Bold         SGRType = 1  // Bold
	SGR_Underline    SGRType = 4  // Underline
	SGR_Blinking     SGRType = 5  // Blinking
	SGR_Negative     SGRType = 7  // Negative image
	SGR_Invisible    SGRType = 8  // Invisible image
	SGR_BoldOff      SGRType = 22 // Bold off
	SGR_UnderlineOff SGRType = 24 // Underline off
	SGR_BlinkingOff  SGRType = 25 // Blinking off
	SGR_NegativeOff  SGRType = 27 // Negative image off
	SGR_InvisibleOff SGRType = 28 // Invisible image off
)

// SGR - Select Graphic Rendition, other data may follow
func SGR(n ...SGRType) string {
	if len(n) < 1 {
		return ""
	}
	s := CSI
	s += fmt.Sprintf("%d", n[0])
	for i := 1; i < len(n); i++ {
		s += fmt.Sprintf(";%d", n[i])
	}
	s += "m"
	return s
}

const (
	SBox_Horiz = 0x2501 //single line
	SBox_Vert  = 0x2503
	SBox_UL    = 0x250F
	SBox_UR    = 0x2513
	SBox_LL    = 0x2517
	SBox_LR    = 0x251B

	DBox_Horiz = 0x2550 //double line
	DBox_Vert  = 0x2551
	DBox_UL    = 0x2554
	DBox_UR    = 0x2557
	DBox_LL    = 0x255A
	DBox_LR    = 0x255D

	BBox_Horiz = 0x2509 //broken line
	BBox_Vert  = 0x250B
	BBox_UL    = 0x250F
	BBox_UR    = 0x2513
	BBox_LL    = 0x2517
	BBox_LR    = 0x251B

	HBox_Horiz = 0x2501 //horiz only
	HBox_Vert  = 0x20
	HBox_UL    = 0x20
	HBox_UR    = 0x20
	HBox_LL    = 0x20
	HBox_LR    = 0x20
)
