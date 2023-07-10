package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate wit jira",
	Long: `Authorize CLI with Jira, so that it can fetch data from Jira.

Note: It's recommended to use API token instead of authenticating with oauth from CLI,
as Atlassian currently doesn't support PKCE verification for oauth flow.`,
	ValidArgs: []string{"login", "logout"},
	Args:      cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "login":
			fmt.Println("auth login")
		case "logout":
			fmt.Println("auth logout")
		default:
			return fmt.Errorf("Invalid argument %s", args[0])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
