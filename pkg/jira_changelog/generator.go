package jira_changelog

import (
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/jira"
	"golang.org/x/exp/slog"
)

type Changelog struct {
	JiraConfig jira.Config
	fromRef    string
	toRef      string
	client     *jira.Client
}

func (c Changelog) Generate() {
	issue, err := c.client.FetchIssue("random-id") // TODO: use correct id
	if err != nil {
		slog.Error("failed while fetching issues from jira", "error", err)
		panic(err)
	}

	slog.Info("Fetched issue", "issue", issue)
}

func NewChangelog(jiraConfig jira.Config, fromRef, toRef string) *Changelog {
	client := jira.NewClient(jiraConfig)
	return &Changelog{
		jiraConfig,
		fromRef,
		toRef,
		client,
	}
}
