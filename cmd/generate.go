package cmd

import (
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog"
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var (
	fromRef string
	toRef   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates changelog",
	Run: func(cmd *cobra.Command, args []string) {
		changelog := jira_changelog.NewGenerator(
			jira.Config{
				BaseUrl:     viper.GetString("base_url"),
				ProjectName: viper.GetString("project_name"),
				ApiToken:    viper.GetString("api_token"),
				User:        viper.GetString("email_id")},
			fromRef,
			toRef,
		)

		slog.Info("Generating changelog", "JiraConfig", changelog.JiraConfig,
			"From", fromRef, "To", toRef)
		changelog.Generate()

		slog.Info("Successfully generated")
	},
}

func init() {
	generateCmd.Flags().StringVar(&fromRef, "from", "", "Git ref to start from")
	generateCmd.Flags().StringVar(&toRef, "to", "main", "Git ref to end at")

	generateCmd.MarkFlagRequired("from")

	rootCmd.AddCommand(generateCmd)
}
