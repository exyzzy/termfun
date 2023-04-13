package termfun

import (
	"fmt"
	"strings"

	"github.com/exyzzy/termfun/format"
)

// TileHandler handles the rendering (output) and keypresses (input) for a specific type of tile
// TileType_ScrollDown renders from top to bottom and breaks lines between words if possible, it scolls up, down.
// TileType_ScrollDownClip renders from top to bottom and only breaks lines on newline in the text, it scrolls up, down, left, right.
// TileType_ScrollDownClipRaw is like TileType_ScrollDownClip but without any scroll key handling.
// TileType_ScrollUp behaves like a typical terminal, rendering bottom to top by line, it scrolls up, down.

type TileType int

const (
	TileType_ScrollDown TileType = iota
	TileType_ScrollDownClip
	TileType_ScrollDownClipRaw
	TileType_ScrollUp
)

type TileHandler struct {
	TileType TileType
	Render   func(*Tile) string
	KeyPress func(*Tile, rune) bool
}

var tileHandler = []*TileHandler{
	{TileType: TileType_ScrollDown, Render: sd_RenderText, KeyPress: sd_KeyPress},
	{TileType: TileType_ScrollDownClip, Render: sdc_RenderText, KeyPress: sdc_KeyPress},
	{TileType: TileType_ScrollDownClipRaw, Render: sdc_RenderText, KeyPress: sdcr_KeyPress},
	{TileType: TileType_ScrollUp, Render: su_RenderText, KeyPress: su_KeyPress}}

// Note that all render handlers must render full text boundary area, clearing as necessary

// renderLines renders top to bottom any lines that line within the tile boundary
// using blank lines as needed
func renderLinesDown(t *Tile, lines []string, origin func(*Tile)) string {
	var str, newLine string
	blankLine := strings.Repeat(" ", t.Width())
	for y := t.bounds.Min.Y; y <= t.bounds.Max.Y; y++ {
		str += CUP(t.bounds.Min.X, y)
		index := y - t.bounds.Min.Y + t.start.Y
		if index >= 0 && index < len(lines) {
			newLine = lines[index]
		} else {
			newLine = blankLine
		}
		str += newLine
	}
	origin(t)
	return str
}

// == TileType_ScrollDown Handler Functions
func sd_RenderText(t *Tile) string {
	tabs := 3
	lines := format.FormatTextBreak(t.buffer.String(), t.Width(), tabs)
	return (renderLinesDown(t, lines, sd_SetCurPosOrigin))
}

func sd_SetCurPosOrigin(t *Tile) {
	t.curPos.X = t.bounds.Min.X
	t.curPos.Y = t.bounds.Min.Y
}

func sd_KeyPress(t *Tile, r rune) bool {
	switch r {
	case KeyUp:
		if t.start.Y > 0 {
			t.start.Y--
			t.setDirty()
		}
	case KeyDown:
		t.start.Y++
		t.setDirty()
	}
	return false // return true from any keypress handler to exit TileTerm
}

// == TileType_ScrollDownClip Handler Functions
func sdc_RenderText(t *Tile) string {
	tabs := 3
	lines := format.FormatTextClipCol(t.buffer.String(), t.Width(), tabs, t.start.X)
	return (renderLinesDown(t, lines, sd_SetCurPosOrigin))

}

func sdc_KeyPress(t *Tile, r rune) bool {
	switch r {
	case KeyUp:
		if t.start.Y > 0 {
			t.start.Y--
			t.setDirty()
		}
	case KeyDown:
		t.start.Y++
		t.setDirty()
	case KeyLeft:
		if t.start.X > 0 {
			t.start.X--
			t.setDirty()
		}
	case KeyRight:
		t.start.X++
		t.setDirty()
	}
	return false
}

// == TileType_ScrollDownClipRaw Handler Functions

// returning true from any KeyCallback will exit TileTerm
type KeyCallback func(r rune) bool

func sdcr_KeyPress(t *Tile, r rune) bool {
	if t.keyCallback != nil {
		return t.keyCallback(r)
	}
	return false
}

// == TileType_ScrollUp Handler Functions
func su_RenderText(t *Tile) string {
	var str string
	blankLine := strings.Repeat(" ", t.Width())
	ss := strings.Split(t.buffer.String(), "\n")
	ss[len(ss)-1] = t.cursor + t.line // add cursor
	curSS := len(ss) - 1
	var sub []string
	curSub := -1
	var newLine string
	for y := t.bounds.Max.Y; y >= t.bounds.Min.Y; y-- {
		str += CUP(t.bounds.Min.X, y)
		if curSub >= 0 {
			newLine = sub[curSub]
			curSub--
		} else {
			if curSS >= 0 {
				sub = wrapStr(ss[curSS], t.Width())
				curSub = len(sub) - 1
				curSS--
				newLine = sub[curSub]
				curSub--
			} else {
				newLine = blankLine
			}
		}
		str += newLine
	}
	su_SetCurPosOrigin(t)
	t.curPos.X += len(t.cursor) + len(t.line)
	return str
}

// wrapStr returns single string s as []strings that are clipped or space padded
func wrapStr(s string, ln int) []string {
	var str string
	var out []string
	if ln == 0 {
		return out
	}
	str = s
	if len(str) == 0 {
		str = " " //always one line min
	}
	for len(str) > 0 {
		if len(str) > ln {
			out = append(out, str[:ln])
			str = str[ln:]
		} else {
			out = append(out, str+strings.Repeat(" ", ln-len(str)))
			str = ""
		}
	}
	return out
}

func su_SetCurPosOrigin(t *Tile) {
	t.curPos.X = t.bounds.Min.X
	t.curPos.Y = t.bounds.Max.Y
}

// returning true from any LineCallback will exit TileTerm
type LineCallback func(string) bool

func su_KeyPress(t *Tile, r rune) bool {
	switch r {
	case KeyEnter:
		if t.lineCallback != nil {
			if t.lineCallback(t.line) {
				return true
			}
			t.historyIndex = -1
			t.history.Add(t.line)

			t.line = ""
		}
	// case KeyLeft, KeyRight:

	case KeyUp:
		entry, ok := t.history.NthPreviousEntry(t.historyIndex + 1)
		if !ok {
			return false
		}
		if t.historyIndex == -1 {
			t.historyPending = t.line
		}
		t.historyIndex++
		t.setLine(entry)

	case KeyDown:
		switch t.historyIndex {
		case -1:
			return false
		case 0:
			t.setLine(t.historyPending)
			t.historyIndex--
		default:
			entry, ok := t.history.NthPreviousEntry(t.historyIndex - 1)
			if ok {
				t.historyIndex--
				t.setLine(entry)
			}
		}

	case KeyBackspace:
		if len(t.line) > 0 {
			t.line = t.line[:len(t.line)-1]
		}

	default:
		t.line += fmt.Sprintf("%c", r)
	}
	return false
}
