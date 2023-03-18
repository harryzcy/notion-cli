package api

import (
	"context"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
	"go.zcy.dev/notion-cli/internal/oauth2"
)

type DatabasePageRecoverInput struct {
	Database string
	PageID   string
}

// RecoverPages recovers a page in a database
func (db database) RecoverPage(input DatabasePageRecoverInput) error {
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

	page, err := client.Page.Get(ctx, notionapi.PageID(input.PageID))
	if err != nil {
		return err
	}

	if page.Parent.Type != "database_id" || page.Parent.DatabaseID != id {
		return fmt.Errorf("page not found")
	}

	if !page.Archived {
		return fmt.Errorf("page is not archived")
	}

	_, err = client.Page.Update(ctx, notionapi.PageID(input.PageID), &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{},
		Archived:   false,
	})
	if err != nil {
		return err
	}

	fmt.Println("page recovered")
	return nil
}
