package database

import (
	"fmt"
	"os"

	"github.com/harryzcy/notion-cli/internal/api"
	"github.com/spf13/cobra"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Args:  cobra.ExactArgs(2),
	Short: "Delete a page from a database",
	Run: func(cmd *cobra.Command, args []string) {
		// delete a page from a database
		input := api.DatabasePageDeleteInput{
			Database: args[0],
			PageID:   args[1],
		}

		err := api.Database.DeletePage(input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
