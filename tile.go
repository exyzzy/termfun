package termfun

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Point struct {
	X, Y int
}

type Rect struct {
	Min, Max Point
}

// LocType is a location within the parent
type LocType int

const (
	Loc_Top LocType = iota
	Loc_Bottom
	Loc_Left
	Loc_Right
)

// Tile contains the state for a Tile
type Tile struct {
	handler      *TileHandler
	name         string          // tile title
	bounds       Rect            // tile bounds
	buffer       strings.Builder // tile text buffer
	cursor       string          // for tile types that support cursors, or ""
	outline      *[6]int         // outline / border character set, or nil for no outline
	curPos       Point           // current position of the tile cursor
	fraction     float32         // fraction of parent this tile uses
	location     LocType         // location within parent for this tile
	parent       *Tile           // parent, or nil if root
	line         string          // current input line for tiles that use it
	dirty        bool            // if true re-render tile
	start        Point           // x,y start of rendering in doc, for scrolling
	keyCallback  KeyCallback
	lineCallback LineCallback
	lock         sync.Mutex

	// stRingBuffer is directly borrowed from golang term
	// history contains previously entered commands so that they can be
	// accessed with the up and down keys.
	history stRingBuffer
	// historyIndex stores the currently accessed history entry, where zero
	// means the immediately previous entry.
	historyIndex int
	// When navigating up and down the history it's possible to return to
	// the incomplete, initial line. That value is stored in
	// historyPending.
	historyPending string
}

// Width returns the Tile's Width
func (t *Tile) Width() int {
	return t.bounds.Max.X - t.bounds.Min.X + 1
}

// Height returns the Tile's Height
func (t *Tile) Height() int {
	return t.bounds.Max.Y - t.bounds.Min.Y + 1
}

// ResetBuffer resets the Tile's string buffer
func (t *Tile) ResetBuffer() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.buffer.Reset()
}

// Cursor returns the current Cursor string for the Tile
func (t *Tile) Cursor() string {
	return t.cursor
}

// Line returns the current Line string for the Tile
func (t *Tile) Line() string {
	return t.line
}

// setLine sets the current Line string in the Tile
func (t *Tile) setLine(line string) {
	t.line = line
	return
}

// SetKeyCallback sets the Key Callback function for TileType_ScrollDownClipRaw
func (t *Tile) SetKeyCallback(c KeyCallback) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.handler.TileType == TileType_ScrollDownClipRaw {
		t.keyCallback = c
		return nil
	} else {
		return errors.New("Handler.TileType does not support Callback")
	}
}

// SetLineCallback sets the Line Callback function for TileType_ScrollUp
func (t *Tile) SetLineCallback(c LineCallback) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.handler.TileType == TileType_ScrollUp {
		t.lineCallback = c
		return nil
	} else {
		return errors.New("Handler.TileType does not support Callback")
	}
}

// Write to support io.Writer interface
func (tile *Tile) Write(buf []byte) (n int, err error) {
	tile.setDirty()
	tile.lock.Lock()
	defer tile.lock.Unlock()
	n, err = tile.buffer.Write(buf)
	return
}

// Print text to a tile buffer
func (tile *Tile) Print(s ...any) {
	tile.setDirty()
	tile.lock.Lock()
	defer tile.lock.Unlock()
	fmt.Fprint(&tile.buffer, s...)
}

// Println text to a tile buffer
func (tile *Tile) Println(s ...any) {
	tile.setDirty()
	tile.lock.Lock()
	defer tile.lock.Unlock()
	fmt.Fprintln(&tile.buffer, s...)
}

// Printf text to a tile buffer
func (tile *Tile) Printf(format string, s ...any) {
	tile.setDirty()
	tile.lock.Lock()
	defer tile.lock.Unlock()
	fmt.Fprintf(&tile.buffer, format, s...)
}

// setDirty sets the tile status to dirty
func (tile *Tile) setDirty() {
	tile.lock.Lock()
	defer tile.lock.Unlock()
	tile.dirty = true
}

//===== Render Helper Functions

// renderOutline renders the tile outline
func (t *Tile) renderOutline(focus bool) string {
	var str string
	if t.outline != nil {
		if focus {
			str += Box(IncRect(t.bounds), *(t.outline), t.name, SGR_Negative, SGR_Bold)
		} else {
			str += Box(IncRect(t.bounds), *(t.outline), t.name)
		}
	}
	return str
}

// Clear clears the tile
func (t *Tile) Clear() string {
	return ClearRect(t.bounds)
}

// HLine is a helper to make a horizontal line using a character
func HLine(x1, x2, y, c int) string {
	return CUP(x1, y) + strings.Repeat(fmt.Sprintf("%c", c), x2-x1)
}

// HLineText is a helper to make a horizontal line with text in the center
// SGRType is applied to the text
func HLineText(x1, x2, y, c int, text string, style ...SGRType) string {
	h := (x2 - x1 - len(text)) / 2
	var str string
	str = CUP(x1, y) + strings.Repeat(fmt.Sprintf("%c", c), h)
	str += SGR(style...) + text + SGR(0)
	str += strings.Repeat(fmt.Sprintf("%c", c), x2-x1-len(text)-h)
	return str
	// return CUP(x1, y) + strings.Repeat(fmt.Sprintf("%c", c), h) + text + strings.Repeat(fmt.Sprintf("%c", c), x2-x1-len(text)-h)
}

// Vline is a helper to make a vertical line using a character
func VLine(x, y1, y2, c int) string {
	return CUP(x, y1) + strings.Repeat(fmt.Sprintf("%c%s%s", c, CUB(1), CUD(1)), y2-y1)
}

// CharAt place a character at a location
func CharAt(x, y, c int) string {
	return CUP(x, y) + fmt.Sprintf("%c", c)
}

// ClearRect returns a clear rectangle of spaces
func ClearRect(r Rect) string {
	if r.Max.X-r.Min.X < 1 || r.Max.Y-r.Min.Y < 1 {
		return ""
	}
	var str string
	spc := strings.Repeat(" ", r.Max.X-r.Min.X+1)
	for j := r.Min.Y; j <= r.Max.Y; j++ {
		str += CUP(r.Min.X, j)
		str += spc
	}
	return str
}

// DecRect returns a rect decreased by 1, to allow for an outline
func DecRect(r Rect) Rect {
	return Rect{Min: Point{X: r.Min.X + 1, Y: r.Min.Y + 1}, Max: Point{X: r.Max.X - 1, Y: r.Max.Y - 1}}
}

// IncRect returns a rect increased by 1
func IncRect(r Rect) Rect {
	return Rect{Min: Point{X: r.Min.X - 1, Y: r.Min.Y - 1}, Max: Point{X: r.Max.X + 1, Y: r.Max.Y + 1}}
}

// box part indices within a char set
const (
	Box_Horiz = iota
	Box_Vert
	Box_UL
	Box_UR
	Box_LL
	Box_LR
)

// box char sets
var SingleBox = [6]int{SBox_Horiz, SBox_Vert, SBox_UL, SBox_UR, SBox_LL, SBox_LR}

var DoubleBox = [6]int{DBox_Horiz, DBox_Vert, DBox_UL, DBox_UR, DBox_LL, DBox_LR}

var BrokenBox = [6]int{BBox_Horiz, BBox_Vert, BBox_UL, BBox_UR, BBox_LL, BBox_LR}

var HorizBox = [6]int{HBox_Horiz, HBox_Vert, HBox_UL, HBox_UR, HBox_LL, HBox_LR}

// Box returns a box using a char set, with optional top title that can be styled
func Box(r Rect, chars [6]int, text string, style ...SGRType) string {
	if r.Max.X-r.Min.X < 2 || r.Max.Y-r.Min.Y < 2 {
		return ""
	}
	var str string
	str = CharAt(r.Min.X, r.Min.Y, chars[Box_UL])
	str += CUU(1)
	if len(text) > (r.Max.X - r.Min.X) {
		str += HLine(r.Min.X+1, r.Max.X, r.Min.Y, chars[Box_Horiz])
	} else {
		str += HLineText(r.Min.X+1, r.Max.X, r.Min.Y, chars[Box_Horiz], text, style...)
	}
	str += CharAt(r.Max.X, r.Min.Y, chars[Box_UR])
	str += VLine(r.Max.X, r.Min.Y+1, r.Max.Y, chars[Box_Vert])
	str += CharAt(r.Max.X, r.Max.Y, chars[Box_LR])
	str += HLine(r.Min.X+1, r.Max.X, r.Max.Y, chars[Box_Horiz])
	str += CharAt(r.Min.X, r.Max.Y, chars[Box_LL])
	str += VLine(r.Min.X, r.Min.Y+1, r.Max.Y, chars[Box_Vert])
	return str
}
