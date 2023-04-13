package termfun

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"golang.org/x/term"
)

// TileTerm contains the state for a TileTerm session
type TileTerm struct {
	width  int           // width of all combined tiles
	height int           // height of all combined tiles
	big    *Tile         // if not nil then this tile is enlarged
	focus  *Tile         // this tile has input focus
	tiles  []*Tile       // all current tiles
	dirty  bool          // if true re-render all tiles
	in     *os.File      // input
	out    *os.File      // ouput
	reader *bufio.Reader // input reader
	lock   sync.Mutex    // protect TileTerm from concurrent processing issues
}

// NewTileTerm returns a new TileTerm session, typically pass in stdin and stdout
func NewTileTerm(in, out *os.File) *TileTerm {
	//make a *bufio.reader
	reader := bufio.NewReader(in)

	return &TileTerm{dirty: true, in: in, reader: reader, out: out}
}

// AddTile adds a new tile to the TileTerm session
// name: Title name of the tile, or ""
// cursor: Cursor suffix for the tile, or ""
// outline: Box character set to use for tile border
// fraction: % of parent tile to take for this child
// location: location within parent tile to use for this child
// parent: which existing tile to use as parent for this child (root is nil)
// handler: which tile handler set to use
func (tTerm *TileTerm) AddTile(name string, cursor string, outline [6]int, fraction float32, location LocType, parent *Tile, handler TileType) (*Tile, error) {
	if len(tTerm.tiles) > 0 && parent == nil {
		return nil, fmt.Errorf("no parent tile")
	}
	tTerm.lock.Lock()
	defer tTerm.lock.Unlock()

	tile := Tile{name: name, cursor: cursor, outline: &outline, fraction: fraction, location: location, parent: parent, handler: tileHandler[handler], historyIndex: -1}
	tTerm.tiles = append(tTerm.tiles, &tile)
	if len(tTerm.tiles) == 1 {
		tTerm.focus = &tile
	}
	return &tile, nil
}

// DeleteTile deletes a tile and all of its children
func (tTerm *TileTerm) DeleteTile(tile *Tile) {
	tTerm.lock.Lock()
	defer tTerm.lock.Unlock()
	if tTerm.focus == tile {
		if tile != tTerm.tiles[0] {
			tTerm.focus = tTerm.tiles[0]
		} else {
			tTerm.focus = nil
		}
	}
	var newTiles []*Tile
	for _, v := range tTerm.tiles {
		if v != tile {
			newTiles = append(newTiles, v)
		}
	}
	tTerm.tiles = newTiles
	tTerm.deleteTileChildren(tile)
}

// TileByIndex returns a tile by its index
func (tTerm *TileTerm) TileByIndex(index int) (tile *Tile) {
	if index >= 0 && index < len(tTerm.tiles) {
		return tTerm.tiles[index]
	} else {
		return nil
	}
}

// String returns the current string of the rendered TileTerm session
func (tTerm *TileTerm) String() string {
	var str string
	w, h, err := term.GetSize(int(tTerm.in.Fd()))
	if err != nil {
		panic(err)
	}

	cw, ch := tTerm.getSize()

	if cw != w || ch != h || tTerm.dirty {
		tTerm.dirtyAllTiles()
		tTerm.setSize(w, h)
		tTerm.layoutTiles()
		str += tTerm.renderOutline()
	}
	str += tTerm.renderText()
	return str
}

// Start begins the TileTerm session
func (tTerm *TileTerm) Start() error {
	w, h, err := term.GetSize(int(tTerm.in.Fd()))
	if err != nil {
		return err
	}
	tTerm.setSize(w, h)

	for {
		// read a key rune from the reader
		key, _, err := ReadKey(tTerm.reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// handle the key
		if tTerm.handleKey(key) {
			return nil
		}
	}
	return nil
}

// Render renders the TileTerm session to Out
// It is a convenience function for String()
func (tTerm *TileTerm) Render() {
	fmt.Fprint(tTerm.out, CUP(1, 1), tTerm.String())
}

// deleteTileChildren deletes the child tiles attached to a parent
func (tTerm *TileTerm) deleteTileChildren(tile *Tile) {
	for _, v := range tTerm.tiles {
		if v.parent == tile {
			tTerm.DeleteTile(v)
		}
	}
}

// setSize sets the size of the TileTerm session
func (tTerm *TileTerm) setSize(width, height int) {
	tTerm.lock.Lock()
	defer tTerm.lock.Unlock()
	tTerm.width, tTerm.height = width, height
}

// getSize returns the size of the TileTerm session
func (tTerm *TileTerm) getSize() (width, height int) {
	return tTerm.width, tTerm.height
}

// Set Dirty sets the session dirty flag which causes rendering of all tiles
func (tTerm *TileTerm) setDirty() {
	tTerm.lock.Lock()
	defer tTerm.lock.Unlock()
	tTerm.dirty = true
}

// nextTile returns the next tile in order of creation, looping back to the first
func (tTerm *TileTerm) nextTile(cur *Tile) (*Tile, error) {
	for i, t := range tTerm.tiles {
		if t == cur {
			return tTerm.tiles[(i+1)%len(tTerm.tiles)], nil
		}
	}
	return nil, errors.New("no tile match")
}

// layoutTiles computes all tile sizes and curpos based on terminal width and height
func (tTerm *TileTerm) layoutTiles() error {
	tTerm.lock.Lock()
	defer tTerm.lock.Unlock()
	tw := tTerm.width
	th := tTerm.height
	for i, w := range tTerm.tiles {
		fraction := w.fraction
		var wr Rect
		if i == 0 {
			wr = Rect{Min: Point{X: 1, Y: 1}, Max: Point{X: tw, Y: th}}
		} else {
			pr := w.parent.bounds
			if tTerm.big != nil {
				if w == tTerm.big {
					fraction = .9
				} else if w == tTerm.big.parent {
					fraction = .9
				} else {
					fraction = .1
				}
			}
			if w.parent.outline != nil {
				pr = IncRect(pr) //grab the outline space
			}
			wr = pr
			switch w.location {
			case Loc_Top:
				dh := int(float32(pr.Max.Y-pr.Min.Y+1)*fraction) + pr.Min.Y
				wr.Max.Y = dh - 1
				pr.Min.Y = dh
			case Loc_Bottom:
				dh := int(float32(pr.Max.Y-pr.Min.Y+1)*(1.-fraction)) + pr.Min.Y
				pr.Max.Y = dh - 1
				wr.Min.Y = dh
			case Loc_Left:
				dw := int(float32(pr.Max.X-pr.Min.X+1)*fraction) + pr.Min.X
				wr.Max.X = dw - 1
				pr.Min.X = dw
			case Loc_Right:
				dw := int(float32(pr.Max.X-pr.Min.X+1)*(1.-fraction)) + pr.Min.X
				pr.Max.X = dw - 1
				wr.Min.X = dw
			}
			if w.parent.outline != nil {
				pr = DecRect(pr)
			}
			w.parent.bounds = pr
		}
		if w.outline != nil {
			wr = DecRect(wr)
		}
		w.bounds = wr
	}
	return nil
}

// renderCursor renders the cursor of the focus tile
func (tTerm *TileTerm) renderCursor() string {
	if tTerm.focus == nil {
		return ""
	}
	return CUP(tTerm.focus.curPos.X, tTerm.focus.curPos.Y)
}

// renderOutline renders the borders of all tiles
func (tTerm *TileTerm) renderOutline() string {
	var str string
	for _, w := range tTerm.tiles {
		str += w.renderOutline(tTerm.focus == w)
	}
	return str
}

// renderText renders the text of any dirty tile
func (tTerm *TileTerm) renderText() string {
	var str string
	for _, t := range tTerm.tiles {
		if t.dirty {
			str += t.handler.Render(t)
		}
	}
	str += tTerm.renderCursor()

	return str
}

// dirtyAllTiles sets all tiles to dirty
func (tTerm *TileTerm) dirtyAllTiles() {
	tTerm.lock.Lock()
	defer tTerm.lock.Unlock()
	for _, t := range tTerm.tiles {
		t.setDirty()
	}
}

// handleKey processes the keys common to all tiles
// continue to call this until it returns true
func (tTerm *TileTerm) handleKey(key rune) bool {
	var err error
	switch key {
	case CtrlT: // tab focus to next tile
		tTerm.focus, err = tTerm.nextTile(tTerm.focus)
		if err != nil {
			panic(err)
		}
		tTerm.setDirty()

	case CtrlU: // enlarge this tile
		if tTerm.big == nil {
			tTerm.big = tTerm.focus
		} else {
			tTerm.big = nil
		}
		tTerm.setDirty()

	case CtrlQ: // quit TileTerm session
		return true

	default: // pass key to tile for handling
		if tTerm.focus.handler.KeyPress(tTerm.focus, key) {
			return true
		}
	}
	tTerm.Render()
	return false
}
