package jira_changelog

import (
	"context"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/messages"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
)

type Generator struct {
	JiraConfig *jira.Context
	fromRef    string
	toRef      string
	repoURL    string
	client     jira.Client
	usePR      bool
}

func NewGenerator(client jira.Client, usePR bool, fromRef, toRef, repoURL string) *Generator {
	g := &Generator{
		fromRef: fromRef,
		toRef:   toRef,
		repoURL: repoURL,
		client:  client,
		usePR:   usePR,
	}

	return g
}

func (c *Generator) Generate(ctx context.Context) *Changelog {
	populator, err := messages.NewCommitOrPRPopualtor(c.usePR, c.fromRef, c.toRef, c.repoURL)
	panicIfErr(err)

	commits, err := populator.Populate(ctx)
	panicIfErr(err)

	issues, err := c.fetchJiraIssues(commits)
	panicIfErr(err)

	issuesByEpic := lo.GroupBy(issues, func(issue jira.Issue) string { return issue.Epic() })

	slog.Debug("Total epics", "count", len(issuesByEpic))

	return NewChangelog(c.fromRef, c.toRef, c.repoURL, issuesByEpic)
}

func (c *Generator) fetchJiraIssues(commits []messages.Message) ([]jira.Issue, error) {
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

func (c *Generator) fetchJiraIssue(commit messages.Message) (jira.Issue, error) {
	issueId := jira.IssueId(commit.Message())
	if issueId == "" {
		slog.Warn("commit message does not contain issue jira id of the project", "commit", commit)
		return jira.NewIssue("", commit.Message(), "done", ""), nil
	}

	issue, err := c.client.FetchIssue(string(issueId))
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
