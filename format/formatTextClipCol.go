package format

import (
	"strings"
	"unicode"
)

// FormatTextClipCol preserves tabs and newlines. 
// Lines longer than width are clipped at width.
// All lines space padded to width.
// Each single visible rune counts toward width,
// beginning at column col.
func FormatTextClipCol(text string, width int, tabSize int, col int) []string {
	text = strings.Replace(text, "\r\n", "\n", -1)
	text = strings.Replace(text, "\r", "\n", -1)
	text = strings.Replace(text, "\t", strings.Repeat(" ", tabSize), -1)
	lines := strings.Split(text, "\n")
	blankLine := strings.Repeat(" ", width)
	var b strings.Builder
	for i, line := range lines {
		var c, j int
		var r rune
		b.Reset()
		for _, r = range line {
			if c >= col && c < col+width {
				b.WriteRune(r)
				j++
			}
			if unicode.IsPrint(r) {
				c++
			}
		}
		if j == 0 {
			lines[i] = blankLine
		} else {
			lines[i] = b.String()
			if c >= col && c < col+width {
				lines[i] += strings.Repeat(" ", col+width-c)
			}
		}
	}
	return lines
}
