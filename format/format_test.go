package format

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

// go test -run TestFormatBreak
func TestFormatBreak(t *testing.T) {
	text := "Loremipsumdolorsitametconsecteturadipiscingelit. \n\n\nMauris        eget purus arcu. Sed quis ornare magna. \n\t- Nulla facilisi\n\t- Praesent in elit\n\t- In mi aliquet suscipit\nSuspendisse\tvel\tenim id metus iaculis pretium. Sed semper pharetra mi a varius. Vestibulum\trutrum\tultricies urna,\tvitae pretium metus ullamcorper vel.\nInteger euismod elit elit, at dictum urna auctor a. Suspendisse vitae est aliquam, euismod enim a, imperdiet ipsum. Suspendisse vel enim id metus iaculis pretium.5Ὂg̀9! ℃ᾭG\n5Ὂg̀9! ℃ᾭG\n<the end>"

	for width := 10; width < 200; width += 23 {
		lines := FormatTextBreak(text, width, 3)
		printCheck(t, lines, width)
	}
}
func printCheck(t *testing.T, lines []string, width int) {
	for _, line := range lines {
		if utf8.RuneCountInString(line) != width {
			fmt.Printf("|%s|%d\n", line, utf8.RuneCountInString(line))
			t.Errorf("Width expected: %d but got: %d", width, utf8.RuneCountInString(line))
		}
	}
}

// go test -run TestFormatClip
func TestFormatClip(t *testing.T) {
	text := "Loremipsumdolorsitametconsecteturadipiscingelit. \n\n\nMauris        eget purus arcu. Sed quis ornare magna. \n\t- Nulla facilisi\n\t- Praesent in elit\n\t- In mi aliquet suscipit\nSuspendisse\tvel\tenim id metus iaculis pretium. Sed semper pharetra mi a varius. Vestibulum\trutrum\tultricies urna,\tvitae pretium metus ullamcorper vel.\nInteger euismod elit elit, at dictum urna auctor a. Suspendisse vitae est aliquam, euismod enim a, imperdiet ipsum. Suspendisse vel enim id metus iaculis pretium.5Ὂg̀9! ℃ᾭG\n5Ὂg̀9! ℃ᾭG\n<the end>"

	col := 0
	for width := 10; width < 200; width += 23 {
		lines := FormatTextClipCol(text, width, 3, col)
		printCheck(t, lines, width)
		col = col + 5
	}
}

// go test -run TestFormatVisual
func TestFormatVisual(t *testing.T) {
	text := "Loremipsumdolorsitametconsecteturadipiscingelit. \n\n\nMauris        eget purus arcu. Sed quis ornare magna. \n\t- Nulla facilisi\n\t- Praesent in elit\n\t- In mi aliquet suscipit\nSuspendisse\tvel\tenim id metus iaculis pretium. Sed semper pharetra mi a varius. Vestibulum\trutrum\tultricies urna,\tvitae pretium metus ullamcorper vel.\nInteger euismod elit elit, at dictum urna auctor a. Suspendisse vitae est aliquam, euismod enim a, imperdiet ipsum. Suspendisse vel enim id metus iaculis pretium.5Ὂg̀9! ℃ᾭG\n5Ὂg̀9! ℃ᾭG\n<the end>"

	width := 40

	fmt.Println("=== FormatTextBreak ===")
	lines := FormatTextBreak(text, width, 3)
	for _, line := range lines {
		fmt.Printf("|%s|%d\n", line, utf8.RuneCountInString(line))
	}
	fmt.Println("=== FormatTextClipCol ===")
	lines = FormatTextClipCol(text, width, 3, 4)
	for _, line := range lines {
		fmt.Printf("|%s|%d\n", line, utf8.RuneCountInString(line))
	}
}
