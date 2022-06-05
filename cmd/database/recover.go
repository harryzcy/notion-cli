package database

import (
	"fmt"
	"os"

	"github.com/harryzcy/notion-cli/internal/api"
	"github.com/spf13/cobra"
)

// RecoverCmd represents the recover command
var RecoverCmd = &cobra.Command{
	Use:   "recover",
	Args:  cobra.ExactArgs(2),
	Short: "Recover a page from a database",
	Run: func(cmd *cobra.Command, args []string) {
		// recover a page from a database
		input := api.DatabasePageRecoverInput{
			Database: args[0],
			PageID:   args[1],
		}

		err := api.Database.RecoverPage(input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
