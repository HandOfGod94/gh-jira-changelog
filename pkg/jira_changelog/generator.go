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
	JiraConfig jira.Config
	fromRef    string
	toRef      string
	client     jira.Client
}

func panicIfErr(err error, args ...interface{}) {
	slog.Error(err.Error(), args...)
	if err != nil {
		panic(err)
	}
}

func (c Generator) Generate(ctx context.Context) *Changelog {
	gitOutput, err := git.ExecGitLog(ctx, c.fromRef, c.toRef)
	panicIfErr(fmt.Errorf("failed to execute git command. %w", err))

	commits, err := gitOutput.Commits()
	panicIfErr(fmt.Errorf("failed to parse git output. %w", err))

	return c.changelogFromCommits(commits)
}

func (c Generator) changelogFromCommits(commits []git.Commit) *Changelog {
	slog.Debug("Total commit messages", "count", len(commits))

	jiraIssueIds := lo.Map(commits, func(commit git.Commit, index int) jira.JiraIssueId {
		return jira.IssueId(c.JiraConfig.ProjectName, commit.Message)
	})
	jiraIssueIds = lo.Filter(jiraIssueIds, func(jiraIssueId jira.JiraIssueId, index int) bool { return jiraIssueId != "" })
	jiraIssueIds = lo.Uniq(jiraIssueIds)
	slog.Debug("Total jira issues ids", "count", len(jiraIssueIds))

	issues := lo.Map(jiraIssueIds, func(jiraIssueId jira.JiraIssueId, index int) jira.Issue {
		issue, err := c.client.FetchIssue(string(jiraIssueId))
		panicIfErr(fmt.Errorf("failed to fetch issue. %w", err), "issue", jiraIssueId)

		slog.Debug("fetched issue", "issue", issue)
		return issue
	})
	slog.Debug("Total issues", "count", len(issues))

	issuesByEpic := lo.GroupBy(issues, func(issue jira.Issue) string { return issue.Epic() })
	return &Changelog{Changes: issuesByEpic}
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
