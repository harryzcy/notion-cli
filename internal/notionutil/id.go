package notionutil

import "regexp"

func IsNotionID(s string) bool {
	matched, _ := regexp.MatchString("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", s)
	return matched
}
