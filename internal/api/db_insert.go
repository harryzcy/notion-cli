package api

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/jomei/notionapi"
	"go.zcy.dev/notion-cli/internal/oauth2"
)

type DatabaseInsertInput struct {
	Database   string
	Properties map[string]string
	Icon       string
	Cover      string
}

func (db database) Insert(input DatabaseInsertInput) error {
	token, err := oauth2.GetToken()
	if err != nil {
		return fmt.Errorf("token not found, please run `notion auth` first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := notionapi.NewClient(notionapi.Token(token.AccessToken))

	id, err := db.getID(ctx, client, input.Database)
	if err != nil {
		return err
	}

	definedProperties, err := db.retrieveDatabaseProperties(ctx, client, id)
	if err != nil {
		return err
	}

	properties := make(notionapi.Properties)
	for k, v := range input.Properties {
		property, err := parsePropertyFromConfigs(definedProperties, k, v)
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
	page, err := client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: id,
		},
		Icon:       icon,
		Cover:      cover,
		Properties: properties,
	})
	if err != nil {
		return err
	}

	fmt.Println("Page inserted:", page.ID)

	return nil
}

func parsePropertyFromConfigs(defined notionapi.PropertyConfigs, name, value string) (notionapi.Property, error) {
	config, ok := defined[name]
	if !ok {
		return nil, fmt.Errorf("property %s not found", name)
	}

	switch config.(type) {
	case *notionapi.TitlePropertyConfig:
		return notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{Text: &notionapi.Text{Content: value}},
			},
		}, nil

	case *notionapi.RichTextPropertyConfig:
		return notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: value}},
			},
		}, nil

	case *notionapi.NumberPropertyConfig:
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %s", value)
		}
		return notionapi.NumberProperty{
			Number: n,
		}, nil

	case *notionapi.SelectPropertyConfig:
		return notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: value,
			},
		}, nil

	case *notionapi.MultiSelectPropertyConfig:
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

	case *notionapi.DatePropertyConfig:
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

	case *notionapi.FormulaPropertyConfig:
		return nil, fmt.Errorf("property %s not supported", config.GetType())

	case *notionapi.RelationPropertyConfig:
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

	case *notionapi.RollupPropertyConfig:
		return nil, fmt.Errorf("property %s not supported", config.GetType())

	case *notionapi.PeoplePropertyConfig:
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

	case *notionapi.FilesPropertyConfig:
		return nil, fmt.Errorf("property %s not supported", config.GetType())

	case *notionapi.CheckboxPropertyConfig:
		value = strings.ToLower(value)
		boolean := false
		if value == "true" || value == "t" || value == "yes" || value == "y" {
			boolean = true
		}
		return notionapi.CheckboxProperty{
			Checkbox: boolean,
		}, nil

	case *notionapi.URLPropertyConfig:
		return notionapi.URLProperty{
			URL: value,
		}, nil

	case *notionapi.EmailPropertyConfig:
		return notionapi.EmailProperty{
			Email: value,
		}, nil

	case *notionapi.PhoneNumberPropertyConfig:
		return notionapi.PhoneNumberProperty{
			PhoneNumber: value,
		}, nil

	}

	return nil, fmt.Errorf("property %s not supported", config.GetType())
}

func parseIcon(icon string) (*notionapi.Icon, error) {
	isEmoji := false
	if strings.HasPrefix(icon, "emoji,") {
		icon = strings.TrimPrefix(icon, "emoji,")
		isEmoji = true
	}

	emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]`)
	if emojiRegex.MatchString(icon) {
		isEmoji = true
	}

	if isEmoji {
		emoji := notionapi.Emoji(icon)
		return &notionapi.Icon{
			Type:  "emoji",
			Emoji: &emoji,
		}, nil
	}

	icon = strings.TrimPrefix(icon, "external,")
	if icon != "" {
		return &notionapi.Icon{
			Type: "external",
			External: &notionapi.FileObject{
				URL: icon,
			},
		}, nil
	}

	return nil, nil
}

func parseCover(cover string) (*notionapi.Image, error) {
	cover = strings.TrimPrefix(cover, "external,")
	if cover != "" {
		return &notionapi.Image{
			Type: "external",
			External: &notionapi.FileObject{
				URL: cover,
			},
		}, nil
	}

	return nil, nil
}
