package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.zcy.dev/notion-cli/internal/api"
)

// InsertCmd represents the insert command
var InsertCmd = &cobra.Command{
	Use:   "insert",
	Args:  cobra.MinimumNArgs(1),
	Short: "Insert a page into a database",
	Run: func(cmd *cobra.Command, args []string) {
		input := api.DatabaseInsertInput{
			Database:   args[0],
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

		err := api.Database.Insert(input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
