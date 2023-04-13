// Example of mutiple tiled window terminal
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/exyzzy/termfun"
	"github.com/exyzzy/termfun/lorem"
	"golang.org/x/term"
)

func main() {
	// put the terminal in raw mode and save state
	in := os.Stdin
	oldState, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(in.Fd()), oldState)

	fmt.Print(termfun.ED(termfun.EraseAll), "\r\n")

	// make a new TileTerm
	tTerm := termfun.NewTileTerm(in, os.Stdout)

	// add the root tile
	t0, err := tTerm.AddTile(" Main ", "M>", termfun.DoubleBox, 1.0, termfun.Loc_Top, nil, termfun.TileType_ScrollUp)
	if err != nil {
		panic(err)
	}

	// create a lineCallback for the root
	tt := &TermType{tile: t0}
	t0.SetLineCallback(tt.lineHandler)

	// print to the root text buffer
	t0.Println("Hello, t0: root, TileType_ScrollUp, bash lineHandler")

	// add more tiles and text
	err = AddTilesAndText(tTerm, t0)
	if err != nil {
		panic(err)
	}

	// render what we have so far
	tTerm.Render()

	// make a delayed popup counting tile in a go routine
	go counter(tTerm)

	// make a delayed popup life tile in a go routine
	go life(tTerm)

	// start the tTerm
	err = tTerm.Start()
	if err != nil {
		panic(err)
	}
}


type TermType struct {
	tile *termfun.Tile
}

// lineHandler is the root LineCallback
func (tt *TermType) lineHandler(line string) bool {
	tt.tile.Printf("%s%s\n", tt.tile.Cursor(), line)
	out, _ := exec.Command("bash", "-c", line).CombinedOutput()
	tt.tile.Print(string(out))
	if strings.Contains(line, "quit") { //example of how to exit TileTerm from LineCallback
		return true
	}
	return false
}

// life starts the life tile in the go routine
func life(tTerm *termfun.TileTerm) {
	time.Sleep(time.Second * 2)

	t5, err := tTerm.AddTile(" Life ", "", termfun.SingleBox, 0.5, termfun.Loc_Top, tTerm.TileByIndex(3), termfun.TileType_ScrollDownClipRaw)
	if err != nil {
		return
	}
	lt := &LifeType{frameRate: 10}
	err = t5.SetKeyCallback(lt.keyHandler)
	if err != nil {
		panic(err)
	}

	tTerm.Render() //render to set t5 tile size

	//size it to the tile, assume aspect render
	c := termfun.NewCanvas(t5.Width(), t5.Height()*2)
	c.PlotWrap()
	Randomize(c)
	// fmt.Print(termfun.ED(termfun.EraseAll), "\r\n")

	for i := 0; i < 300; i++ {
		t5.ResetBuffer()
		t5.Print(c.StringAspect())
		c = NextFrame(c)
		time.Sleep(time.Second / time.Duration(lt.frameRate))
		tTerm.Render()
	}
	tTerm.DeleteTile(t5)
	tTerm.Render()
}

type LifeType struct {
	frameRate int
}

// keyHandler is the KeyCallback for the life tile
// keydown to reduce framerate, keyup to increase framerate
// q quits TileTerm when life has focus
func (lt *LifeType) keyHandler(k rune) bool {
	if k == termfun.KeyUp {
		if lt.frameRate < 50 {
			lt.frameRate *= 5
		}
	}
	if k == termfun.KeyDown {
		if lt.frameRate > 2 {
			lt.frameRate /= 5
		}
	}
	if k == 'q' { //example of how to exit TileTerm from KeyCallback
		return true 
	}
	return false
}

// AddTilesAndText adds 3 new tiles and some random text, and instructions
func AddTilesAndText(tTerm *termfun.TileTerm, t0 *termfun.Tile) error {

	// make some tiles in the tTerm
	t1, err := tTerm.AddTile(" Text1 ", "", termfun.SingleBox, 0.3, termfun.Loc_Top, t0, termfun.TileType_ScrollDown)
	if err != nil {
		return err
	}
	t2, err := tTerm.AddTile(" Text2 ", "", termfun.HorizBox, 0.3, termfun.Loc_Right, t1, termfun.TileType_ScrollDownClip)
	if err != nil {
		return err
	}
	t3, err := tTerm.AddTile(" Text3 ", "", termfun.SingleBox, 0.3, termfun.Loc_Left, t0, termfun.TileType_ScrollDown)
	if err != nil {
		return err
	}

	// print some text to the tile buffers
	t1.Println("Hello, t1: TileType_ScrollDown")
	t2.Println("Hello, t2: TileType_ScrollDownClip")
	t3.Println("Hello, t3, (tab)\tTileType_ScrollDown")

	t1.Println()
	instructions(t2)
	t2.Println("\nTest of long line:")
	t2.Println("\tAnotherVeryveryveryLongveryveryveryveryEversolongveryveryveryveryveryveryveryveryveryVeryVeryveryveryveryveryveryveryvery, very, very, very, very, very, very, very, very, very, very...long line")
	t3.Println("Another ever so very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very, very ...long line")
	t1.Println(lorem.GenerateLorem(500))
	t3.Println(lorem.GenerateLorem(100))

	return nil
}

// instructions provides instruction text for the tiles
func instructions(t *termfun.Tile) {
	t.Println("What can I do?")
	t.Println("For any/all tile:")
	t.Println("\t- Ctrl-T to cycle focus to next window")
	t.Println("\t- Ctrl-U to make this window big (toggle)")
	t.Println("\t- Ctrl-Q to quit (exit demo)")
	t.Println("For Life:")
	t.Println("\t- Up Arrow to increase frame rate")
	t.Println("\t- Down Arrow to decrease frame rate")
	t.Println("\t- q to quit the whole demo")
	t.Println("For TileType_ScrollDown:")
	t.Println("\t- Up Arrow to scroll up one line")
	t.Println("\t- Down Arrow to scroll down one line")
	t.Println("For TileType_ScrollDownClip:")
	t.Println("\t- Up Arrow to scroll up one line")
	t.Println("\t- Down Arrow to scroll down one line")
	t.Println("\t- Left Arrow to scroll left one column")
	t.Println("\t- Right Arrow to scroll right one column")
	t.Println("For TileType_ScrollUp:")
	t.Println("\t- Enter to send line to bash and print result")
	t.Println("\t- Type 'ls<enter>' for instance")
}

// counter starts the counter tile in the go routine
// print line to tile every second and render
func counter(tTerm *termfun.TileTerm) {
	// new tile pops after 5 seconds
	time.Sleep(time.Second * 5)
	t4, err := tTerm.AddTile(" Text4: Counting tile, a longer title ", "", termfun.SingleBox, 0.3, termfun.Loc_Top, tTerm.TileByIndex(0), termfun.TileType_ScrollDown)
	if err != nil {
		return
	}
	tTerm.Render()
	for i := 0; i <= 20; i++ {
		t4.Println(i)
		tTerm.Render()
		time.Sleep(time.Second * 1)
	}
	tTerm.DeleteTile(t4)
	tTerm.Render()
}

// Life helpers

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
