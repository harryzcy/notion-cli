package api

import "github.com/jomei/notionapi"

func parseRichTextList(rich []notionapi.RichText) string {
	var result string
	for _, r := range rich {
		result += r.PlainText
	}
	return result
}
