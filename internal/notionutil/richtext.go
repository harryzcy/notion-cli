package notionutil

import "github.com/jomei/notionapi"

func ParseRichTextList(rich []notionapi.RichText) string {
	var result string
	for _, r := range rich {
		result += r.PlainText
	}
	return result
}
