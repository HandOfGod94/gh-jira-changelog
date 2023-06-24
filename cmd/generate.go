package cmd

import (
	jira "github.com/handofgod94/jira_changelog/pkg/jira_changelog"
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
		changelog := jira.Changelog{
			JiraConfig: jira.JiraConfig{ProjectUrl: projectUrl, ProjectName: projectName, ApiToken: apiToken},
			GitConfig:  jira.GitConfig{FromRef: fromRef, ToRef: toRef},
		}
		slog.Info("Generating changelog", "Jira Config", changelog.JiraConfig, "Git Config", changelog.GitConfig)
		changelog.Generate()
	},
}

func init() {
	generateCmd.Flags().StringVar(&fromRef, "from", "", "Git ref to start from")
	generateCmd.Flags().StringVar(&toRef, "to", "main", "Git ref to end at")

	generateCmd.MarkFlagRequired("from")

	rootCmd.AddCommand(generateCmd)
}
