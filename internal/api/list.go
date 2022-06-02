package api

import (
	"context"
	"fmt"
	"time"

	"github.com/harryzcy/notion-cli/internal/notionutil"
	"github.com/harryzcy/notion-cli/internal/oauth2"
	"github.com/harryzcy/notion-cli/internal/print"
	"github.com/jomei/notionapi"
)

var (
	ErrMultipleDatabasesFound = fmt.Errorf("multiple databases found")
)

func ListDatabases() error {
	token, err := oauth2.GetToken()
	if err != nil {
		return fmt.Errorf("token not found, please run `notion auth` first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := notionapi.NewClient(notionapi.Token(token.AccessToken))

	print.PrintDatabaseHeader()

	cursor, err := listDatabasePage(ctx, client, "")
	if err != nil {
		return err
	}
	for cursor != "" {
		cursor, err = listDatabasePage(ctx, client, cursor)
		if err != nil {
			return err
		}
	}

	return nil
}

func listDatabasePage(ctx context.Context, client *notionapi.Client, cursor notionapi.Cursor) (notionapi.Cursor, error) {
	res, err := client.Search.Do(ctx, &notionapi.SearchRequest{
		Filter: map[string]string{
			"property": "object",
			"value":    "database",
		},
		PageSize:    100,
		StartCursor: cursor,
	})
	if err != nil {
		return "", err
	}

	for _, result := range res.Results {
		db := result.(*notionapi.Database)

		print.PrintDatabaseEntry(db)
	}
	return res.NextCursor, nil
}

func ListPagesInDatabase(database string) error {
	token, err := oauth2.GetToken()
	if err != nil {
		return fmt.Errorf("token not found, please run `notion auth` first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := notionapi.NewClient(notionapi.Token(token.AccessToken))

	var id notionapi.DatabaseID
	if notionutil.IsNotionID(database) {
		id = notionapi.DatabaseID(database)
	} else {
		id, err = GetDatabaseIDByName(ctx, client, database)
		if err != nil {
			if err == ErrMultipleDatabasesFound {
				err = nil
			}
			return err
		}
	}

	err = queryDatabase(ctx, client, id)

	return err
}

func GetDatabaseIDByName(ctx context.Context, client *notionapi.Client, name string) (notionapi.DatabaseID, error) {
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

func queryDatabase(ctx context.Context, client *notionapi.Client, id notionapi.DatabaseID) error {
	properties, err := retrieveDatabaseProperties(ctx, client, id)
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

func retrieveDatabaseProperties(ctx context.Context,
	client *notionapi.Client, id notionapi.DatabaseID,
) (notionapi.PropertyConfigs, error) {
	res, err := client.Database.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return res.Properties, nil
}
