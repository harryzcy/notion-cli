package print

import (
	"strings"

	"golang.org/x/text/width"
)

func Padding(text string, width int) string {
	occupied := GetWidthUTF8String(text)
	return text + strings.Repeat(" ", width-occupied)
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
