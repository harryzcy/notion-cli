package api

import (
	"testing"
	"time"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/assert"
)

type InvalidPropertyConfig struct {
}

func (p InvalidPropertyConfig) GetType() notionapi.PropertyConfigType {
	return "invalid"
}

func TestParseProperty(t *testing.T) {
	defined := notionapi.PropertyConfigs{
		"title":        &notionapi.TitlePropertyConfig{},
		"rich_text":    &notionapi.RichTextPropertyConfig{},
		"number":       &notionapi.NumberPropertyConfig{},
		"select":       &notionapi.SelectPropertyConfig{},
		"multi_select": &notionapi.MultiSelectPropertyConfig{},
		"date":         &notionapi.DatePropertyConfig{},
		"formula":      &notionapi.FormulaPropertyConfig{},
		"relation":     &notionapi.RelationPropertyConfig{},
		"rollup":       &notionapi.RollupPropertyConfig{},
		"people":       &notionapi.PeoplePropertyConfig{},
		"files":        &notionapi.FilesPropertyConfig{},
		"checkbox":     &notionapi.CheckboxPropertyConfig{},
		"url":          &notionapi.URLPropertyConfig{},
		"email":        &notionapi.EmailPropertyConfig{},
		"phone_number": &notionapi.PhoneNumberPropertyConfig{},
		"invalid":      &InvalidPropertyConfig{},
	}

	property, err := parsePropertyFromConfigs(defined, "title", "hello")
	assert.Nil(t, err)
	assert.Len(t, property.(notionapi.TitleProperty).Title, 1)
	assert.Equal(t, "hello", property.(notionapi.TitleProperty).Title[0].Text.Content)

	property, err = parsePropertyFromConfigs(defined, "rich_text", "hello")
	assert.Nil(t, err)
	assert.Len(t, property.(notionapi.RichTextProperty).RichText, 1)
	assert.Equal(t, "hello", property.(notionapi.RichTextProperty).RichText[0].Text.Content)

	property, err = parsePropertyFromConfigs(defined, "number", "1")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, property.(notionapi.NumberProperty).Number)
	property, err = parsePropertyFromConfigs(defined, "number", "x")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid number x", err.Error())
	assert.Empty(t, property)

	property, err = parsePropertyFromConfigs(defined, "select", "hello")
	assert.Nil(t, err)
	assert.Equal(t, "hello", property.(notionapi.SelectProperty).Select.Name)

	property, err = parsePropertyFromConfigs(defined, "multi_select", "hello,world")
	assert.Nil(t, err)
	assert.Len(t, property.(notionapi.MultiSelectProperty).MultiSelect, 2)
	assert.Equal(t, "hello", property.(notionapi.MultiSelectProperty).MultiSelect[0].Name)
	assert.Equal(t, "world", property.(notionapi.MultiSelectProperty).MultiSelect[1].Name)

	property, err = parsePropertyFromConfigs(defined, "date", "2020-01-01")
	assert.Nil(t, err)
	assert.Equal(t, notionapi.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)), *property.(notionapi.DateProperty).Date.Start)
	property, err = parsePropertyFromConfigs(defined, "date", "x")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid date x", err.Error())
	assert.Empty(t, property)

	property, err = parsePropertyFromConfigs(defined, "formula", "1+1")
	assert.NotNil(t, err)
	assert.Empty(t, property)

	property, err = parsePropertyFromConfigs(defined, "relation", "id")
	assert.Nil(t, err)
	assert.Len(t, property.(notionapi.RelationProperty).Relation, 1)
	assert.EqualValues(t, "id", property.(notionapi.RelationProperty).Relation[0].ID)

	property, err = parsePropertyFromConfigs(defined, "rollup", "")
	assert.NotNil(t, err)
	assert.Empty(t, property)

	property, err = parsePropertyFromConfigs(defined, "people", "foo,bar")
	assert.Nil(t, err)
	assert.Len(t, property.(notionapi.PeopleProperty).People, 2)
	assert.EqualValues(t, "foo", property.(notionapi.PeopleProperty).People[0].ID)
	assert.EqualValues(t, "bar", property.(notionapi.PeopleProperty).People[1].ID)

	property, err = parsePropertyFromConfigs(defined, "files", "hello")
	assert.NotNil(t, err)
	assert.Empty(t, property)

	property, err = parsePropertyFromConfigs(defined, "checkbox", "true")
	assert.Nil(t, err)
	assert.Equal(t, true, property.(notionapi.CheckboxProperty).Checkbox)
	property, err = parsePropertyFromConfigs(defined, "checkbox", "t")
	assert.Nil(t, err)
	assert.Equal(t, true, property.(notionapi.CheckboxProperty).Checkbox)
	property, err = parsePropertyFromConfigs(defined, "checkbox", "yes")
	assert.Nil(t, err)
	assert.Equal(t, true, property.(notionapi.CheckboxProperty).Checkbox)
	property, err = parsePropertyFromConfigs(defined, "checkbox", "y")
	assert.Nil(t, err)
	assert.Equal(t, true, property.(notionapi.CheckboxProperty).Checkbox)
	property, err = parsePropertyFromConfigs(defined, "checkbox", "")
	assert.Nil(t, err)
	assert.Equal(t, false, property.(notionapi.CheckboxProperty).Checkbox)

	property, err = parsePropertyFromConfigs(defined, "url", "hello")
	assert.Nil(t, err)
	assert.Equal(t, "hello", property.(notionapi.URLProperty).URL)

	property, err = parsePropertyFromConfigs(defined, "email", "hello")
	assert.Nil(t, err)
	assert.Equal(t, "hello", property.(notionapi.EmailProperty).Email)

	property, err = parsePropertyFromConfigs(defined, "phone_number", "hello")
	assert.Nil(t, err)
	assert.Equal(t, "hello", property.(notionapi.PhoneNumberProperty).PhoneNumber)

	property, err = parsePropertyFromConfigs(defined, "invalid", "hello")
	assert.NotNil(t, err)
	assert.Empty(t, property)

	property, err = parsePropertyFromConfigs(defined, "unknown", "hello")
	assert.NotNil(t, err)
	assert.Empty(t, property)
}

func emoji(s string) *notionapi.Emoji {
	e := notionapi.Emoji(s)
	return &e
}

func TestParseIcon(t *testing.T) {
	tests := []struct {
		icon        string
		expected    notionapi.Icon
		expectedErr error
	}{
		{
			icon: "😂",
			expected: notionapi.Icon{
				Type:  "emoji",
				Emoji: emoji("😂"),
			},
		},
		{
			icon: "emoji,😂",
			expected: notionapi.Icon{
				Type:  "emoji",
				Emoji: emoji("😂"),
			},
		},
		{
			icon: "https://example.com/emoji.png",
			expected: notionapi.Icon{
				Type: "external",
				External: &notionapi.FileObject{
					URL: "https://example.com/emoji.png",
				},
			},
		},
		{
			icon: "external,https://example.com/emoji.png",
			expected: notionapi.Icon{
				Type: "external",
				External: &notionapi.FileObject{
					URL: "https://example.com/emoji.png",
				},
			},
		},
		{
			icon:     "",
			expected: notionapi.Icon{},
		},
	}

	for _, test := range tests {
		t.Run(test.icon, func(t *testing.T) {
			icon, err := parseIcon(test.icon)
			assert.Equal(t, test.expected, icon)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestParseCover(t *testing.T) {
	tests := []struct {
		cover       string
		expected    notionapi.FileObject
		expectedErr error
	}{
		{
			cover: "https://example.com/cover.png",
			expected: notionapi.FileObject{
				URL: "https://example.com/cover.png",
			},
		},
		{
			cover: "external,https://example.com/cover.png",
			expected: notionapi.FileObject{
				URL: "https://example.com/cover.png",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.cover, func(t *testing.T) {
			cover, err := parseCover(test.cover)
			assert.Equal(t, test.expected, cover)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
