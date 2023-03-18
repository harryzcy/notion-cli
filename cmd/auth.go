package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.zcy.dev/notion-cli/internal/oauth2"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Args:  cobra.NoArgs,
	Short: "Authenticate with Notion",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")

		if clientID == "" {
			clientID = os.Getenv("NOTION_CLIENT_ID")
		}
		if clientSecret == "" {
			clientSecret = os.Getenv("NOTION_CLIENT_SECRET")
		}

		if clientID == "" || clientSecret == "" {
			return fmt.Errorf("client-id and client-secret are required")
		}

		return oauth2.Flow(clientID, clientSecret)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.Flags().StringP("client-id", "i", "", "Notion client ID")
	authCmd.Flags().StringP("client-secret", "s", "", "Notion client secret")
}
