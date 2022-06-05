package cmd

import (
	"github.com/spf13/cobra"

	"github.com/harryzcy/notion-cli/cmd/database"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Handle databases",
}

func init() {
	rootCmd.AddCommand(databaseCmd)

	databaseCmd.AddCommand(database.LsCmd)
	databaseCmd.AddCommand(database.InsertCmd)
	databaseCmd.AddCommand(database.TrashCmd)
	databaseCmd.AddCommand(database.RecoverCmd)
}
