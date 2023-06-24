package jira_changelog

import (
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/jira"
	"golang.org/x/exp/slog"
)

type Generator struct {
	JiraConfig jira.Config
	fromRef    string
	toRef      string
	client     *jira.Client
}

func (c Generator) Generate() {
	issue, err := c.client.FetchIssue("random-id") // TODO: use correct id
	if err != nil {
		slog.Error("failed while fetching issues from jira", "error", err)
		panic(err)
	}

	slog.Info("Fetched issue", "issue", issue)
}

func NewGenerator(jiraConfig jira.Config, fromRef, toRef string) *Generator {
	client := jira.NewClient(jiraConfig)
	return &Generator{
		jiraConfig,
		fromRef,
		toRef,
		client,
	}
}
