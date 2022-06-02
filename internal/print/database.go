package print

import (
	"fmt"
	"strings"
	"time"

	"github.com/harryzcy/notion-cli/internal/notionutil"
	"github.com/jomei/notionapi"
)

const (
	idPadding           = 37
	databasePadding     = 20
	databasePagePadding = 15
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

func PrintDatabasePageHeader(properties notionapi.PropertyConfigs) []string {
	fmt.Print(TruncateOrPad("ID", idPadding))
	orderedNames := make([]string, len(properties))

	i := 0
	for name := range properties {
		fmt.Print(TruncateOrPad(name, databasePagePadding))
		orderedNames[i] = name
		i++

		fmt.Print(" ")
	}
	fmt.Println()

	return orderedNames
}

func PrintDatabasePageEntry(page notionapi.Page, propertyNames []string) {
	fmt.Print(Padding(page.ID.String(), idPadding))
	for _, name := range propertyNames {
		value := ""

		switch v := page.Properties[name].(type) {
		case *notionapi.TitleProperty:
			value = notionutil.ParseRichTextList(v.Title)

		case *notionapi.RichTextProperty:
			value = notionutil.ParseRichTextList(v.RichText)

		case *notionapi.TextProperty:
			value = notionutil.ParseRichTextList(v.Text)

		case *notionapi.NumberProperty:
			value = fmt.Sprintf("%f", v.Number)

		case *notionapi.SelectProperty:
			value = v.Select.Name

		case *notionapi.MultiSelectProperty:
			for _, option := range v.MultiSelect {
				value += option.Name + ", "
			}
			value = strings.TrimSuffix(value, ", ")

		case *notionapi.DateProperty:
			value = formateDateObject(v.Date)

		case *notionapi.FormulaProperty:
			switch v.Formula.Type {
			case notionapi.FormulaTypeString:
				value = v.Formula.String
			case notionapi.FormulaTypeNumber:
				value = fmt.Sprintf("%f", v.Formula.Number)
			case notionapi.FormulaTypeBoolean:
				value = fmt.Sprintf("%t", v.Formula.Boolean)
			case notionapi.FormulaTypeDate:
				value = formateDateObject(v.Formula.Date)
			}

		case *notionapi.RelationProperty:
			for _, relation := range v.Relation {
				value += relation.ID.String() + ", "
			}
			value = strings.TrimSuffix(value, ", ")

		case *notionapi.RollupProperty:
			// TODO: support rollup
			value = "rollup property"

		case *notionapi.PeopleProperty:
			for _, person := range v.People {
				value += person.Name + ", "
			}
			value = strings.TrimSuffix(value, ", ")

		case *notionapi.FilesProperty:
			value = fmt.Sprintf("%d files", len(v.Files))

		case *notionapi.CheckboxProperty:
			value = fmt.Sprintf("%t", v.Checkbox)

		case *notionapi.URLProperty:
			value = v.URL

		case notionapi.EmailProperty:
			value = v.Email

		case notionapi.PhoneNumberProperty:
			value = v.PhoneNumber

		case *notionapi.CreatedTimeProperty:
			value = v.CreatedTime.Format("2006-01-02")

		case notionapi.CreatedByProperty:
			value = v.CreatedBy.Name

		case *notionapi.LastEditedTimeProperty:
			value = v.LastEditedTime.Format("2006-01-02")

		case notionapi.LastEditedByProperty:
			value = v.LastEditedBy.Name

		default:
			value = string(page.Properties[name].GetType()) + "unsupported"
		}

		fmt.Print(TruncateOrPad(value, databasePagePadding) + " ")
	}
	fmt.Println()
}

func formateDateObject(date *notionapi.DateObject) string {
	if date == nil {
		return ""
	}

	value := ""
	if date.Start != nil {
		value += time.Time(*date.Start).Format("2006-01-02")
	}

	if date.End != nil {
		if value != "" {
			value += " - "
		}
		value += time.Time(*date.End).Format("2006-01-02")
	}

	return value
}
