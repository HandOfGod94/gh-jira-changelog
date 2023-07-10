/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate wit jira",
	Long: `Authorize CLI with Jira, so that it can fetch data from Jira.

Note: It's recommended to use API token instead of authenticating with oauth from CLI,
as Atlassian currently doesn't support PKCE verification for oauth flow.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auth called")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
