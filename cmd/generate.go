package cmd

import (
	"context"
	"fmt"
	"io"
	"net/url"
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
	DefaultTimeout = 5 * time.Second
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates changelog",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		apiToken := viper.GetString("api_token")
		emailID := viper.GetString("email_id")
		if apiToken != "" && emailID == "" {
			return fmt.Errorf("valid email_id is required with api_token config")
		}

		_, err := url.Parse(viper.GetString("base_url"))
		if err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

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

	generateCmd.PersistentFlags().StringP("base_url", "u", "", "base url where jira is hosted")
	generateCmd.PersistentFlags().String("email_id", "", "email id of the user")
	generateCmd.PersistentFlags().StringP("api_token", "t", "", "API token for jira")
	generateCmd.PersistentFlags().StringP("project_name", "p", "", "Project name in jira. usually the acronym")
	generateCmd.PersistentFlags().StringP("log_level", "v", "error", "log level. options: debug, info, warn, error")

	generateCmd.MarkFlagRequired("from")
	generateCmd.MarkPersistentFlagRequired("base_url")
	generateCmd.MarkPersistentFlagRequired("project_name")

	rootCmd.AddCommand(generateCmd)
	viper.BindPFlags(generateCmd.PersistentFlags())
}
