package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gh-jira-changelog",
	Short: "Changelog generator using jira issues",
	Long: `Most of our changelog tools solely focus on commits. While the orgs usually use jira to track issues.
When generating changelog why not combine both commits and jira issues to generate a changelog.

This can also work as a plugin for "gh" cli
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var slogLevel slog.Level
		switch viper.GetString("log_level") {
		case "debug":
			slogLevel = slog.LevelDebug
		case "info":
			slogLevel = slog.LevelInfo
		case "warn":
			slogLevel = slog.LevelWarn
		case "error":
			slogLevel = slog.LevelError
		default:
			slogLevel = slog.LevelError
		}

		appLogger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slogLevel})
		slog.SetDefault(slog.New(appLogger))

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.jira_changelog.yaml)")
	rootCmd.PersistentFlags().StringP("base_url", "u", "", "base url where jira is hosted")
	rootCmd.PersistentFlags().String("email_id", "", "email id of the user")
	rootCmd.PersistentFlags().StringP("api_token", "t", "", "API token for jira")
	rootCmd.PersistentFlags().StringP("project_name", "p", "", "Project name in jira. usually the acronym")
	rootCmd.PersistentFlags().StringP("log_level", "v", "error", "log level. options: debug, info, warn, error")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(cwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".jira_changelog")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
