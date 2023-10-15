package cmd

import (
	"context"
	"fmt"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate wit jira",
	Long: `Authorize CLI with Jira, so that it can fetch data from Jira.

Note: It's recommended to use API token instead of authenticating with oauth from CLI,
as Atlassian currently doesn't support PKCE verification for oauth flow.`,
	ValidArgs: []string{"login", "logout"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "login":
			a := jira.NewAuthenticator()
			if err := a.Login(context.Background()); err != nil {
				return err
			}
		case "logout":
			panic("To be implemented")
		default:
			return fmt.Errorf("invalid argument %s", args[0])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
