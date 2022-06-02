package print

import (
	"strings"

	"golang.org/x/text/width"
)

func TruncateOrPad(s string, maxWidth int) string {
	width := GetWidthUTF8String(s)
	if width <= maxWidth {
		return s + strings.Repeat(" ", maxWidth-width)
	}
	return s[:maxWidth]
}

func Padding(text string, width int) string {
	occupied := GetWidthUTF8String(text)
	if width > occupied {
		return text + strings.Repeat(" ", width-occupied)
	}
	return text
}

func GetWidthUTF8String(s string) int {
	size := 0
	for _, runeValue := range s {
		p := width.LookupRune(runeValue)
		if p.Kind() == width.EastAsianWide {
			size += 2
			continue
		}
		if p.Kind() == width.EastAsianNarrow {
			size += 1
			continue
		}
	}
	return size
}
