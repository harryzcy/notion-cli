package database

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/harryzcy/notion-cli/internal/api"
)

// LsCmd represents the ls command
var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List contents",
	Run: func(cmd *cobra.Command, args []string) {
		err := api.ListDatabases()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
