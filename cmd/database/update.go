package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.zcy.dev/notion-cli/internal/api"
)

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update database pageID [property=value...]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Update a page in a database",
	Run: func(cmd *cobra.Command, args []string) {
		input := api.DatabasePageUpdateInput{
			Database:   args[0],
			PageID:     args[1],
			Properties: make(map[string]string),
		}

		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			name := parts[0]
			value := ""
			if len(parts) > 1 {
				value = parts[1]
			}

			if name == "icon" {
				input.Icon = value
			} else if name == "cover" {
				input.Cover = value
			} else {
				input.Properties[name] = value
			}
		}

		err := api.Database.UpdatePage(input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
