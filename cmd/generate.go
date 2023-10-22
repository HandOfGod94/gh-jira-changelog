package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var (
	fromRef        string
	toRef          string
	writeTo        string
	DefaultTimeout = 10 * time.Second
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates changelog",
	Example: `#using as a standalone tool
gh-jira-changelog generate \
	--base_url="<you-atlassian-url>" \
	--from="<git-ref>" \
	--to="<git-ref>" \
	--api_token="<jira-api-token>" \
	--email_id="jira-email-id"
	--repo_url="https://github.com/<org>/<repo>"

# using config file
# all the jira config such as (base_url, api_token, email_id) can be stored in a config file
gh-jira-changelog generate --config="<path-to-config-file>.yaml" --from=<git-ref> --to=<git-ref>

# using env variables
# all the jira config such as (base_url, api_token, email_id) can be provided by env variables
BASE_URL=<you-atlassian-url> API_TOKEN=<jira-api-token> gh-jira-changelog generate --from=<git-ref> --to=<git-ref>

# generating changelog between 2 git tags
gh-jira-changelog generate --config="<path-to-config-file>.yaml" --from="v0.1.0" --to="v0.2.0"


# Using it as GH plugin
# assuming jira plugin installed
gh jira-changelog generate --config="<path-to-config-file>.yaml" --from="v0.1.0" --to="v0.2.0"`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		apiToken := viper.GetString("api_token")
		emailID := viper.GetString("email_id")
		if apiToken != "" && emailID == "" {
			return fmt.Errorf("valid email_id is required with api_token config")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		changelog := jira_changelog.NewGenerator(
			jira.NewContext(jira.Options{
				jira.BaseURL:  viper.GetString("base_url"),
				jira.ApiToken: viper.GetString("api_token"),
				jira.User:     viper.GetString("email_id"),
			}),
			fromRef,
			toRef,
			viper.GetString("repo_url"),
		)

		slog.Info("Generating changelog", "JiraConfig", changelog.JiraConfig,
			"From", fromRef, "To", toRef, "repoURL", viper.GetString("repo_url"))
		changelog.Generate(ctx).Render(writer(writeTo))
	},
}

func writer(writeTo string) io.Writer {
	switch writeTo {
	case "/dev/stdout":
		return os.Stdout
	case "/dev/stderr":
		return os.Stderr
	default:
		file, err := os.Create(writeTo)
		if err != nil {
			slog.Error("error creating output file", "error", err)
			panic(err)
		}
		return file

	}
}

func init() {
	generateCmd.Flags().StringVar(&fromRef, "from", "", "Git ref to start from")
	generateCmd.Flags().StringVar(&toRef, "to", "main", "Git ref to end at")
	generateCmd.Flags().StringVar(&writeTo, "write_to", "/dev/stdout", "File stream to write the changelog")

	generateCmd.MarkFlagRequired("from")

	rootCmd.AddCommand(generateCmd)
}
