package print

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPadding(t *testing.T) {
	tests := []struct {
		text     string
		width    int
		expected string
	}{
		{
			text:     "",
			width:    0,
			expected: "",
		},
		{
			text:     "abc",
			width:    5,
			expected: "abc  ",
		},
		{
			text:     "你好",
			width:    5,
			expected: "你好 ",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.expected, Padding(test.text, test.width))
		})
	}
}

func TestGetWidthUTF8String(t *testing.T) {
	tests := []struct {
		text  string
		width int
	}{
		{
			text:  "abc",
			width: 3,
		},
		{
			text:  "世界",
			width: 4,
		},
		{
			text:  "あいう",
			width: 6,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			width := GetWidthUTF8String(test.text)
			assert.Equal(t, test.width, width)
		})
	}
}
