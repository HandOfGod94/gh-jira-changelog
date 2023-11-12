package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	config "github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira/config_service"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira/oauth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:       "auth",
	Short:     "Authenticate wit jira",
	Long:      `Authorize CLI with Jira using oauth, so that it can fetch data from Jira.`,
	ValidArgs: []string{"login", "logout"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "login":
			a := oauth.NewAuthenticator()
			if err := a.Login(context.Background()); err != nil {
				fmt.Fprintln(os.Stderr, color.RedString("Login attempt failed. Please try again"))
				return err
			}

			fmt.Println(color.GreenString("Login successful"))
		case "logout":
			err := config.Clear()
			if err != nil {
				return err
			}
			fmt.Print(color.CyanString("Logged out successfully."))
			return nil
		default:
			return fmt.Errorf("invalid argument %s", args[0])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
