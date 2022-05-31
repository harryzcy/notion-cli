package cmd

import (
	"github.com/harryzcy/notion-cli/internal/oauth2"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Args:  cobra.NoArgs,
	Short: "Authenticate with Notion",
	RunE: func(cmd *cobra.Command, args []string) error {
		return oauth2.Flow()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
