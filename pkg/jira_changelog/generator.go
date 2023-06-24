package jira_changelog

import (
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/jira"
	"golang.org/x/exp/slog"
)

type Changelog struct {
	JiraConfig jira.Config
	FromRef    string
	ToRef      string
}

func (c Changelog) Generate() {

	slog.Info("Fetching issues from jira")
	client := jira.NewClient(c.JiraConfig)
	issue, err := client.FetchIssue("random-id") // TODO: use correct id
	if err != nil {
		slog.Error("failed while fetching issues from jira", "error", err)
		panic(err)
	}

	slog.Info("Fetched issue", "issue", issue)
}
