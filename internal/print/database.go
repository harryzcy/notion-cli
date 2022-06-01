package print

import (
	"fmt"

	"github.com/harryzcy/notion-cli/internal/notionutil"
	"github.com/jomei/notionapi"
)

func PrintDatabaseHeader() {
	fmt.Printf("%-37s %-20s %-20s %-20s\n",
		"ID",
		"Name",
		"Created At",
		"Edited At",
	)
}

func PrintDatabaseEntry(db *notionapi.Database) {
	title := notionutil.ParseRichTextList(db.Title)
	fmt.Printf("%-37s %s %-20s %-20s\n",
		db.ID,
		Padding(title, 20),
		db.CreatedTime.Format("2006-01-02 15:04:05"),
		db.LastEditedTime.Format("2006-01-02 15:04:05"),
	)
}
