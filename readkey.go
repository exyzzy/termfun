package termfun

import "bufio"

const (
	CtrlA = 0x01 + iota
	CtrlB
	CtrlC
	CtrlD
	CtrlE
	CtrlF
	CtrlG
	CtrlH
	CtrlI
	CtrlJ
	CtrlK
	CtrlL
	CtrlM
	CtrlN
	CtrlO
	CtrlP
	CtrlQ
	CtrlR
	CtrlS
	CtrlT
	CtrlU
	CtrlV
	CtrlW
	CtrlX
	CtrlY
	CtrlZ
)

const (
	KeyTab       = 9
	KeyEnter     = 13
	KeyEscape    = 27
	KeySpace     = 32
	KeyLBracket   = 91
	KeyBackspace = 127
)

const (
	KeyUnknown = 0xd800 /* UTF-16 surrogate area */ + iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyBackTab
	KeyDel
)

// ReadKey is a drop-in replacement for bufio.ReadRune but returns common keyboard keypresses that are multi-rune
// as a single utf-16 surrogate rune per the consts above.
func ReadKey(reader *bufio.Reader) (r rune, size int, err error) {
	var keybuf [4]rune
	var keylen int
	var bts int

	for {
		c, n, err := reader.ReadRune()
		if err != nil {
			return c, n, err
		}
		bts += n
		keybuf[keylen] = c
		keylen++

		if keylen == 1 && keybuf[0] != KeyEscape {
			return keybuf[0], keylen, nil
		}
		if keylen == 2 && keybuf[1] != KeyLBracket {
			err := reader.UnreadRune() // not a CSI
			if err != nil {
				return c, n, err
			}
			return keybuf[0], keylen - n, nil
		}
		if keylen == 3 {
			switch keybuf[2] {
			case 65:
				return KeyUp, bts, nil
			case 66:
				return KeyDown, bts, nil
			case 67:
				return KeyRight, bts, nil
			case 68:
				return KeyLeft, bts, nil
			case 90:
				return KeyBackTab, bts, nil
			case 51:
				// nothing
			default:
				return KeyUnknown, bts, nil
			}
		}
		if keylen == 4 {
			if keybuf[2] == 51 && keybuf[3] == 126 {
				return KeyDel, bts, nil
			} else {
				return KeyUnknown, bts, nil
			}
		}
	}
}
