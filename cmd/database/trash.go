package database

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.zcy.dev/notion-cli/internal/api"
)

// TrashCmd represents the trash command
var TrashCmd = &cobra.Command{
	Use:   "trash",
	Args:  cobra.ExactArgs(2),
	Short: "Trash a page from a database",
	Run: func(cmd *cobra.Command, args []string) {
		// trash a page from a database
		input := api.DatabasePageTrashInput{
			Database: args[0],
			PageID:   args[1],
		}

		err := api.Database.TrashPage(input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
