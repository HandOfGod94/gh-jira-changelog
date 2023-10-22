package jira_changelog

import (
	"context"
	"fmt"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/git"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
)

type Generator struct {
	JiraConfig *jira.Context
	fromRef    string
	toRef      string
	repoURL    string
	client     jira.Client
}

func NewGenerator(jiraCtx *jira.Context, fromRef, toRef, repoURL string) *Generator {
	client := jira.NewClient(jiraCtx)
	g := &Generator{
		JiraConfig: jiraCtx,
		fromRef:    fromRef,
		toRef:      toRef,
		repoURL:    repoURL,
		client:     client,
	}

	return g
}

func (c *Generator) Generate(ctx context.Context) *Changelog {
	commits, err := git.NewCommitPopulator(c.fromRef, c.toRef).Commits(ctx)
	panicIfErr(err)

	issues, err := c.fetchJiraIssues(commits)
	panicIfErr(err)

	issuesByEpic := lo.GroupBy(issues, func(issue jira.Issue) string { return issue.Epic() })

	slog.Debug("Total epics", "count", len(issuesByEpic))

	return NewChangelog(c.fromRef, c.toRef, c.repoURL, issuesByEpic)
}

func (c *Generator) fetchJiraIssues(commits []git.Commit) ([]jira.Issue, error) {
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

func (c *Generator) fetchJiraIssue(commit git.Commit) (jira.Issue, error) {
	issueId := jira.IssueId(commit.Message)
	if issueId == "" {
		slog.Warn("commit message does not contain issue jira id of the project", "commit", commit)
		return jira.NewIssue("", fmt.Sprintf("%s (%s)", commit.Message, commit.Sha), "done", ""), nil
	}

	issue, err := c.client.FetchIssue(string(issueId))
	if err != nil {
		slog.Warn("failed to fetch jira issue", "commit", commit)
		return jira.NewIssue("", fmt.Sprintf("%s (%s)", commit.Message, commit.Sha), "done", ""), nil
	}
	return issue, nil
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
