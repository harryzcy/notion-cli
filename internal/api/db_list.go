package api

import (
	"context"
	"fmt"
	"time"

	"github.com/harryzcy/notion-cli/internal/oauth2"
	"github.com/harryzcy/notion-cli/internal/print"
	"github.com/jomei/notionapi"
)

var (
	ErrMultipleDatabasesFound = fmt.Errorf("multiple databases found")
)

func (db database) List() error {
	token, err := oauth2.GetToken()
	if err != nil {
		return fmt.Errorf("token not found, please run `notion auth` first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := notionapi.NewClient(notionapi.Token(token.AccessToken))

	print.PrintDatabaseHeader()

	cursor, err := db.getDBPage(ctx, client, "")
	if err != nil {
		return err
	}
	for cursor != "" {
		cursor, err = db.getDBPage(ctx, client, cursor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db database) getDBPage(ctx context.Context, client *notionapi.Client, cursor notionapi.Cursor) (notionapi.Cursor, error) {
	res, err := client.Search.Do(ctx, &notionapi.SearchRequest{
		Filter: notionapi.SearchFilter{
			Value:    "database",
			Property: "object",
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

func (db database) ListPages(database string) error {
	token, err := oauth2.GetToken()
	if err != nil {
		return fmt.Errorf("token not found, please run `notion auth` first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := notionapi.NewClient(notionapi.Token(token.AccessToken))

	id, err := db.getID(ctx, client, database)
	if err != nil {
		return err
	}

	err = db.queryDatabase(ctx, client, id)
	return err
}
