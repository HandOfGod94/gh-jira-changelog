package cmd

import (
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog"
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/jira"
	"github.com/spf13/cobra"
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
		changelog := jira_changelog.Changelog{
			JiraConfig: jira.Config{ProjectUrl: projectUrl, ProjectName: projectName, ApiToken: apiToken},
			FromRef:    fromRef,
			ToRef:      toRef,
		}

		slog.Info("Generating changelog", "JiraConfig", changelog.JiraConfig,
			"From", changelog.FromRef, "To", changelog.ToRef)
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
