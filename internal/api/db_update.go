package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/harryzcy/notion-cli/internal/oauth2"
	"github.com/jomei/notionapi"
)

type DatabasePageUpdateInput struct {
	Database   string
	PageID     string
	Properties map[string]string
	Icon       string
	Cover      string
}

func (db database) UpdatePage(input DatabasePageUpdateInput) error {
	token, err := oauth2.GetToken()
	if err != nil {
		return fmt.Errorf("token not found, please run `notion auth` first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := notionapi.NewClient(notionapi.Token(token.AccessToken))

	dbID, err := db.getID(ctx, client, input.Database)
	if err != nil {
		return err
	}

	page, err := client.Page.Get(ctx, notionapi.PageID(input.PageID))
	if err != nil {
		return err
	}
	if page.Parent.Type != "database_id" || page.Parent.DatabaseID != dbID {
		return fmt.Errorf("page not found")
	}

	properties := make(notionapi.Properties)
	for k, v := range input.Properties {
		property, err := parseProperty(page.Properties, k, v)
		if err != nil {
			return err
		}
		properties[k] = property
	}

	icon, err := parseIcon(input.Icon)
	if err != nil {
		return err
	}
	cover, err := parseCover(input.Cover)
	if err != nil {
		return err
	}

	_, err = client.Page.Update(ctx, notionapi.PageID(input.PageID), &notionapi.PageUpdateRequest{
		Properties: properties,
		Icon:       icon,
		Cover:      cover,
	})

	return err
}

func parseProperty(defined notionapi.Properties, name, value string) (notionapi.Property, error) {
	config, ok := defined[name]
	if !ok {
		return nil, fmt.Errorf("property %s not found", name)
	}

	switch config.GetType() {
	case "title":
		return notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{Text: &notionapi.Text{Content: value}},
			},
		}, nil

	case "rich_text":
		return notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: value}},
			},
		}, nil

	case "number":
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %s", value)
		}
		return notionapi.NumberProperty{
			Number: n,
		}, nil

	case "select":
		return notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: value,
			},
		}, nil

	case "multi_select":
		values := strings.Split(value, ",")
		p := notionapi.MultiSelectProperty{
			MultiSelect: make([]notionapi.Option, len(values)),
		}
		for i, v := range values {
			p.MultiSelect[i] = notionapi.Option{
				Name: v,
			}
		}
		return p, nil

	case "date":
		t, err := dateparse.ParseAny(value)
		if err != nil {
			return nil, fmt.Errorf("invalid date %s", value)
		}
		notionDate := notionapi.Date(t)
		return notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &notionDate,
			},
		}, nil

	case "formula":
		return nil, fmt.Errorf("property %s not supported", config.GetType())

	case "relation":
		values := strings.Split(value, ",")
		property := notionapi.RelationProperty{
			Relation: make([]notionapi.Relation, len(values)),
		}
		for i, v := range values {
			property.Relation[i] = notionapi.Relation{
				ID: notionapi.PageID(v),
			}
		}
		return property, nil

	case "rollup":
		return nil, fmt.Errorf("property %s not supported", config.GetType())

	case "people":
		values := strings.Split(value, ",")
		property := notionapi.PeopleProperty{
			People: make([]notionapi.User, len(values)),
		}
		for i, v := range values {
			property.People[i] = notionapi.User{
				ID: notionapi.UserID(v),
			}
		}
		return property, nil

	case "files":
		return nil, fmt.Errorf("property %s not supported", config.GetType())

	case "checkbox":
		value = strings.ToLower(value)
		boolean := false
		if value == "true" || value == "t" || value == "yes" || value == "y" {
			boolean = true
		}
		return notionapi.CheckboxProperty{
			Checkbox: boolean,
		}, nil

	case "url":
		return notionapi.URLProperty{
			URL: value,
		}, nil

	case "email":
		return notionapi.EmailProperty{
			Email: value,
		}, nil

	case "phone_number":
		return notionapi.PhoneNumberProperty{
			PhoneNumber: value,
		}, nil

	}

	return nil, fmt.Errorf("property %s not supported", config.GetType())
}
