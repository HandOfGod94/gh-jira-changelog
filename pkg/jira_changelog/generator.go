package jira_changelog

import (
	"context"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/git"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
)

type Generator struct {
	JiraConfig jira.Config
	fromRef    string
	toRef      string
	client     jira.Client
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (c Generator) Generate(ctx context.Context) *Changelog {
	gitOutput, err := git.ExecGitLog(ctx, c.fromRef, c.toRef)
	panicIfErr(err)

	commits, err := gitOutput.Commits()
	panicIfErr(err)

	changelog, err := c.changelogFromCommits(commits)
	panicIfErr(err)

	return changelog
}

func (c Generator) changelogFromCommits(commits []git.Commit) (*Changelog, error) {
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

	jiraIssues = lo.Uniq(jiraIssues)
	slog.Debug("Total jira issues ids", "count", len(jiraIssues))

	issuesByEpic := lo.GroupBy(jiraIssues, func(issue jira.Issue) string { return issue.Epic() })
	return &Changelog{Changes: issuesByEpic}, nil
}

func (c Generator) fetchJiraIssue(commit git.Commit) (jira.Issue, error) {
	issueId := jira.IssueId(c.JiraConfig.ProjectName, commit.Message)
	if issueId == "" {
		slog.Warn("commit message does not contain issue jira id of the project", "commit", commit, "project", c.JiraConfig.ProjectName)
		return jira.NewIssue("", commit.Message, "done", ""), nil
	}

	issue, err := c.client.FetchIssue(string(issueId))
	if err != nil {
		return jira.Issue{}, err
	}
	return issue, nil
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
