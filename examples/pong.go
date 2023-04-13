// Pong style game - press left or right arrow keys to move paddle. 
// Paddle has 5 velocities: -2, -1, 0, 1, 2. 
// Shift between them with left and right arrow keys.
// Try to hit the ball with the paddle, any miss is recorded.
// Ball x velocity depends on where it hits the paddle.
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/exyzzy/termfun"
	"golang.org/x/term"
)

type Game struct {
	canvas     *termfun.Canvas
	paddleSize int
	paddleX    int
	paddleV    int
	key        chan rune
	ballX      int
	ballY      int
	ballVx     int
	ballVy     int
	frameRate  int
	misses     int
	start      time.Time
}

const Width = 110 //pixels
const Height = 50 //pixels
const Paddle = 12 //pixels
const FrameRate = 20

func main() {
	in := os.Stdin
	oldState, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(in.Fd()), oldState)

	// hold game state
	g := &Game{
		canvas:     termfun.NewCanvas(Width, Height),
		paddleSize: Paddle,
		paddleX:    Width/2 - Paddle/2,
		paddleV:    1,
		key:        make(chan rune),
		ballX:      Width / 2,
		ballY:      Height - 3,
		ballVx:     -1,
		ballVy:     -1,
		frameRate:  FrameRate,
		misses:     0,
		start:      time.Now(),
	}

	fmt.Print(termfun.ED(termfun.EraseAll), "\r\n")

	// animate in separate process
	go Animate(g)

	reader := bufio.NewReader(in)
	var k rune
	fmt.Print("Ponglike (<- and -> to shift paddle speed, q to quit)\r\n")
	for {
		k, _, err = termfun.ReadKey(reader)
		if err != nil {
			panic(err)
		}
		switch k {
		case termfun.KeyLeft, termfun.KeyRight:
			g.key <- k
		case termfun.CtrlC, 'q':
			return
		}

	}
}

func Animate(g *Game) {
	for {
		DrawFrame(g)
		fmt.Print(termfun.CUP(0, 0), g.canvas.StringDenseBorder(), "\r\n")
		fmt.Printf("Misses: %4d  Elapsed: %4d\r\n", g.misses, int(time.Since(g.start).Seconds()))
		HandleEvents(g)
		time.Sleep(time.Second / time.Duration(g.frameRate))
	}
}

func HandleEvents(g *Game) {
	if g.ballY >= g.canvas.Height()-1 {
		if g.ballX >= g.paddleX && g.ballX <= g.paddleX+g.paddleSize {
			g.ballVx = -2 + int(float32(g.ballX-g.paddleX)/float32(g.paddleSize)*5.0)
			g.ballY -= g.ballVy
			g.ballVy = -g.ballVy
		} else {
			Miss(g)
		}
		return
	}

	// shift paddle velocity
	select {
	case key := <-g.key:
		switch key {
		case termfun.KeyLeft:
			if g.paddleV > -2 {
				g.paddleV--
			}
		case termfun.KeyRight:
			if g.paddleV < 2 {
				g.paddleV++
			}
		}
	default: // don't block
	}
	g.paddleX += g.paddleV
	if g.paddleX < 0 || g.paddleX+g.paddleSize >= g.canvas.Width() {
		g.paddleX -= g.paddleV
		g.paddleV = -g.paddleV
	}
	g.ballX += g.ballVx
	g.ballY += g.ballVy

	if g.ballX < 0 || g.ballX >= g.canvas.Width() {
		g.ballX -= g.ballVx
		g.ballVx = -g.ballVx
	}
	if g.ballY < 0 || g.ballY >= g.canvas.Height() {
		g.ballY -= g.ballVy
		g.ballVy = -g.ballVy
	}
}

// increment miss count, small delay, serve again
func Miss(g *Game) {
	g.paddleX = g.canvas.Width()/2 - Paddle/2
	g.paddleV = 1
	g.misses++
	g.ballX = g.canvas.Width() / 2
	g.ballY = g.canvas.Height() - 3
	g.ballVx = -1
	g.ballVy = -1
	time.Sleep(time.Second)
}

// just draw the ball and the paddle, frame on string render
func DrawFrame(g *Game) {
	g.canvas.Clear()
	g.canvas.Plot(g.ballX, g.ballY)
	g.canvas.Line(g.paddleX, g.canvas.Height()-1, g.paddleX+g.paddleSize, g.canvas.Height()-1)
	return
}
