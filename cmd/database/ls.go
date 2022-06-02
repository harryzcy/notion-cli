package database

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/harryzcy/notion-cli/internal/api"
)

// LsCmd represents the ls command
var LsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "List databases or pages in a database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := api.Database.List()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		err := api.Database.ListPages(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
