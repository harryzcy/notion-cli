package api

import (
	"context"
	"fmt"

	"github.com/harryzcy/notion-cli/internal/notionutil"
	"github.com/harryzcy/notion-cli/internal/print"
	"github.com/jomei/notionapi"
)

type DatabaseAPI interface {
	List() error
	ListPages(database string) error
	Insert(input DatabaseInsertInput) error
	TrashPage(input DatabasePageTrashInput) error
}

type database struct{}

var Database DatabaseAPI = database{}

// getID returns the database ID from the given either ID or name.
func (db database) getID(ctx context.Context, client *notionapi.Client, database string) (notionapi.DatabaseID, error) {
	if notionutil.IsNotionID(database) {
		return notionapi.DatabaseID(database), nil
	}

	id, err := db.getDatabaseIDByName(ctx, client, database)
	if err != nil {
		if err == ErrMultipleDatabasesFound {
			err = nil
		}
		return "", err
	}

	return id, nil
}

func (db database) getDatabaseIDByName(ctx context.Context, client *notionapi.Client, name string) (notionapi.DatabaseID, error) {
	res, err := client.Search.Do(ctx, &notionapi.SearchRequest{
		Query: name,
		Filter: map[string]string{
			"property": "object",
			"value":    "database",
		},
	})
	if err != nil {
		return "", err
	}

	matched := make([]*notionapi.Database, 0)

	for _, result := range res.Results {
		db := result.(*notionapi.Database)
		if notionutil.ParseRichTextList(db.Title) == name {
			matched = append(matched, db)
		}
	}

	if len(matched) == 0 {
		return "", fmt.Errorf("database not found: %s", name)
	}

	if len(matched) == 1 {
		return notionapi.DatabaseID(matched[0].ID), nil
	}

	fmt.Println("multiple databases found, please specify by id")
	print.PrintDatabaseHeader()
	for _, db := range matched {
		print.PrintDatabaseEntry(db)
	}

	return "", ErrMultipleDatabasesFound
}

func (db database) queryDatabase(ctx context.Context, client *notionapi.Client, id notionapi.DatabaseID) error {
	properties, err := db.retrieveDatabaseProperties(ctx, client, id)
	if err != nil {
		return err
	}

	orderedNames := print.PrintDatabasePageHeader(properties)

	hasMore := true
	nextCursor := notionapi.Cursor("")
	for hasMore {
		res, err := client.Database.Query(ctx, id, &notionapi.DatabaseQueryRequest{
			StartCursor: nextCursor,
		})
		if err != nil {
			return err
		}

		hasMore = res.HasMore
		nextCursor = res.NextCursor
		for _, page := range res.Results {
			print.PrintDatabasePageEntry(page, orderedNames)
		}
	}

	return nil
}

func (db database) retrieveDatabaseProperties(ctx context.Context,
	client *notionapi.Client, id notionapi.DatabaseID,
) (notionapi.PropertyConfigs, error) {
	res, err := client.Database.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return res.Properties, nil
}
