package jira_changelog

import (
	"context"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/messages"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
)

type Generator struct {
	FromRef   string
	ToRef     string
	RepoURL   string
	Client    jira.Client
	Populator messages.Populator
}

func (c *Generator) Generate(ctx context.Context) *Changelog {
	commits, err := c.Populator.Populate(ctx)
	panicIfErr(err)

	issues, err := c.fetchJiraIssues(commits)
	panicIfErr(err)

	issuesByEpic := lo.GroupBy(issues, func(issue jira.Issue) string { return issue.Epic() })

	slog.Debug("Total epics", "count", len(issuesByEpic))

	return NewChangelog(c.FromRef, c.ToRef, c.RepoURL, issuesByEpic)
}

func (c *Generator) fetchJiraIssues(commits []messages.Messager) ([]jira.Issue, error) {
	slog.Debug("Total commit messages", "count", len(commits))

	jiraIssues := make([]jira.Issue, 0)
	for _, commit := range commits {
		issue, err := c.fetchJiraIssue(commit)
		if err != nil {
			slog.Error("error fetching jira issue", "error", err, "commit", commit)
			return nil, err
		}

		slog.Debug("fetched issue", "issue", issue)
		jiraIssues = append(jiraIssues, issue)
	}
	return lo.Uniq(jiraIssues), nil
}

func (c *Generator) fetchJiraIssue(commit messages.Messager) (jira.Issue, error) {
	issueId := jira.IssueId(commit.Message())
	if issueId == "" {
		slog.Warn("commit message does not contain issue jira id of the project", "commit", commit)
		return jira.NewIssue("", commit.Message(), "done", ""), nil
	}

	issue, err := c.Client.FetchIssue(string(issueId))
	if err != nil {
		slog.Warn("failed to fetch jira issue", "commit", commit)
		return jira.NewIssue("", commit.Message(), "done", ""), nil
	}
	return issue, nil
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
