package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	projectUrl  string
	projectName string
	apiToken    string
)

var rootCmd = &cobra.Command{
	Use:   "jira-changelog",
	Short: "Changelog generator using jira issues",
	Long: `Most of our changelog tools solely focus on commits. While the orgs usually use jira to track issues.
When generating changelog why not combine both commits and jira issues to generate a changelog.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		_, err := url.Parse(projectUrl)
		if err != nil {
			return err
		}
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
	rootCmd.PersistentFlags().StringVarP(&projectUrl, "project_url", "u", "", "jira project url")
	rootCmd.PersistentFlags().StringVarP(&apiToken, "api_token", "t", "", "API token for jira")
	rootCmd.PersistentFlags().StringVarP(&projectName, "project_name", "p", "", "Project name in jira. usually the acronym")

	rootCmd.MarkPersistentFlagRequired("project_url")
	rootCmd.MarkPersistentFlagRequired("api_token")
	rootCmd.MarkPersistentFlagRequired("project_name")
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
