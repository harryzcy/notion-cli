package api

import (
	"context"
	"fmt"
	"time"

	"github.com/harryzcy/notion-cli/internal/oauth2"
	"github.com/jomei/notionapi"
)

func ListDatabases() error {
	token, err := oauth2.GetToken()
	if err != nil {
		return fmt.Errorf("token not found, please run `notion auth` first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := notionapi.NewClient(notionapi.Token(token.AccessToken))

	res, err := client.Search.Do(ctx, &notionapi.SearchRequest{
		Filter: map[string]string{
			"property": "object",
			"value":    "database",
		},
	})
	if err != nil {
		return err
	}
	for _, result := range res.Results {
		database := result.(*notionapi.Database)
		fmt.Printf("%s %s\n", database.ID, database.Title[0].PlainText)
	}

	return nil
}
