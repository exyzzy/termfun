# TermFun

TermFun is a package to do some fun things in your Go Terminal.

## ReadKey

```
func ReadKey(reader *bufio.Reader) (r rune, size int, err error)
```
ReadKey is a drop-in replacement for bufio.ReadRune but returns common keyboard keypresses that are multi-rune as a single utf-16 surrogate rune.

Example: ./examples/readkey.go
Example: ./examples/pong.go

```
	KeyUnknown
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyBackTab
	KeyDel
```


## CSI Codes

A collection of CSI functions to control the raw terminal display. All functions return strings that can be printed to the raw display. 

See: https://en.wikipedia.org/wiki/ANSI_escape_code#CSI_(Control_Sequence_Introducer)_sequences

Example: ./examples/box.go

Example: ./examples/tile.go

```
// CUU - Cursor Up
func CUU(n int) string 

// CUD - Cursor Down
func CUD(n int) string

// CUF - Cursor Forward
func CUF(n int) string

// CUB - Cursor Back
func CUB(n int) string

// CNL - Cursor Next Line
func CNL(n int) string 

// CPL - Cursor Previous Line
func CPL(n int) string

// CHA - Cursor Horizontal Absolute
func CHA(n int) string

// CUP - Cursor Position
func CUP(m, n int) string 

// ED - Erase in Display
func ED(n EraseType) string

// SU - Scroll Up
func SU(n int) string 

// HVP - Horizontal Vertical Position
func HVP(m, n int) string

// SGR - Select Graphic Rendition, other data may follow
func SGR(n ...SGRType) string 
```

# Canvas

The Canvas type and methods allow simple 2d graphics with unicode block characters.

Example: ./examples/life.go
Example: ./examples/pong.go

```
// Return a new initialized Canvas
func NewCanvas(w, h int) *Canvas

// Init an existing Canvas, w and h are in pixels
// So Init(80,80) renders an array of 40x40 characters as StringDense
// and an array of 80x40 characters as StringAspect
func (c *Canvas) Init(w, h int)

// Width returns the width of the canvas, in pixels
func (c *Canvas) Width() int  { return c.pwidth }

// Height returns the height of the canvas, in pixels
func (c *Canvas) Height() int { return c.pheight }

// PlotXor sets the plot mode to xor the new pixel on the existing pixel
func (c *Canvas) PlotXor()    { c.xor = true }

// PlotOr sets the plot mode to or the new pixel on the existing pixel (default)
func (c *Canvas) PlotOr()     { c.xor = false }

// Plotwrap sets the plot mode to wrap cordinates on the borders
func (c *Canvas) PlotWrap()   { c.wrap = true }

// Plotwrap sets the plot mode to not wrap cordinates on the borders (default)
func (c *Canvas) PlotUnwrap() { c.wrap = false }

// Plot a pixel on the Canvas
func (c *Canvas) Plot(x, y int)

// Read a location from the Canvas to see if it is set
func (c *Canvas) Read(x, y int) bool

// Clear the Canvas
func (c *Canvas) Clear()

// Render the Canvas to a string as "square" (1x2) pixels
func (c *Canvas) StringAspect() string

// Render the Canvas to a string as "square" (1x2) pixels, with a border
func (c *Canvas) StringAspectBorder() string

// Render the Canvas to a string as dense (2x2) rectangle pixels
func (c *Canvas) StringDense() string

// Render the Canvas to a string as dense (2x2) rectangle pixels, with a border
func (c *Canvas) StringDenseBorder() string

// Line uses Bresenham's algorithm to plot a line
func (c *Canvas) Line(x0, y0, x1, y1 int) 

// Bmp plots a [][]byte to the Canvas
func (c *Canvas) Bmp(x, y int, bmp [][]byte)

```

# TileTerm

The TileTerm type and methods allow rendering multiple tiled display regions. CtrlT will cycle between all the tiles, giving "focus" to each tile in order. CtrlU will make the current focus tile bigger and squish the surrounding tiles. CtrlU again will undo that. The focus tile will act differently on the keypresses depending on its TileType and the specific application.

Example: ./examples/tile.go


### TileType_ScrollDown:
Use this for nicely flow formatting text within the horizontal boundary. The lines try to automatically break at spaces. You can scroll up/down to see more.

KeyUp: scroll up

KeyDown: scroll down

### TileType_ScrollDownClip:
Use this to preserve existing line formatting. You may need to scroll left/right to see all of the clipped region.

KeyUp: scroll up

KeyDown: scroll down

KeyLeft: scroll left

KeyRight: scroll right

### TileType_ScrollDownClipRaw:
Use this for completely custom key handling, you must handle all keys in keyCallBack. It can also be used to disable all keys in a tile.

All keys: call the keyCallBack function

### TileType_ScrollUp:
Use this to scroll up, like a traditional terminal. The return key sends the line to the lineCallBack function for application handling. The up/down keys support a command ring buffer.

KeyEnter: send line to lineCallBack function and add line to ringbuffer.

KeyUp: cycle through ringbuffer and replace line

KeyDown: uncycle through ringbuffer and replace line

KeyBackspace: delete last character from line

### TileTerm API

```
// NewTileTerm returns a new TileTerm session, typically pass in stdin and stdout
func NewTileTerm(in, out *os.File) *TileTerm

// AddTile adds a new tile to the TileTerm session
// name: Title name of the tile, or ""
// cursor: Cursor suffix for the tile, or ""
// outline: Box character set to use for tile border
// fraction: % of parent tile to take for this child
// location: location within parent tile to use for this child
// parent: which existing tile to use as parent for this child (root is nil)
// handler: which tile handler set to use
func (tTerm *TileTerm) AddTile(name string, cursor string, outline [6]int, fraction float32, location LocType, parent *Tile, handler TileType) (*Tile, error)

// DeleteTile deletes a tile and all of its children
func (tTerm *TileTerm) DeleteTile(tile *Tile)

// TileByIndex returns a tile by its index
func (tTerm *TileTerm) TileByIndex(index int) (tile *Tile)

// String returns the current string of the rendered TileTerm session
func (tTerm *TileTerm) String() string 

// Start begins the TileTerm session
func (tTerm *TileTerm) Start() error 

// Render renders the TileTerm session to Out
// It is a convenience function for String()
func (tTerm *TileTerm) Render() 
```

### Tile API

```
// Width returns the Tile's Width
func (t *Tile) Width() int

// Height returns the Tile's Height
func (t *Tile) Height() int

// ResetBuffer resets the Tile's string buffer
func (t *Tile) ResetBuffer()

// Cursor returns the current Cursor string for the Tile
func (t *Tile) Cursor() string 

// Line returns the current Line string for the Tile
func (t *Tile) Line() string 

// SetKeyCallback sets the Key Callback function for TileType_ScrollDownClipRaw
func (t *Tile) SetKeyCallback(c KeyCallback) error

// SetLineCallback sets the Line Callback function for TileType_ScrollUp
func (t *Tile) SetLineCallback(c LineCallback) error

// Write to support io.Writer interface, so you can also, for instance,  fmt.Fprint(tile, "Hello")
func (tile *Tile) Write(buf []byte) (n int, err error)

// Print text to a tile buffer, convenience for Write
func (tile *Tile) Print(s ...any) 

// Println text to a tile buffer, convenience for Write
func (tile *Tile) Println(s ...any) 

// Printf text to a tile buffer, convenience for Write
func (tile *Tile) Printf(format string, s ...any)
```