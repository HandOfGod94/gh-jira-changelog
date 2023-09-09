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
	repoURL    string
	client     jira.Client
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (c Generator) Generate(ctx context.Context) *Changelog {
	gcw := git.NewCommitParseWorkflow(c.fromRef, c.toRef)
	commits, err := gcw.Commits(ctx)
	if err != nil {
		panic(fmt.Errorf("failed at \"%s\" state. %w. State: %+v", gcw.FSM.Current(), err, gcw))
	}

	changes, err := c.changesFromCommits(commits)
	panicIfErr(err)

	slog.Debug("changes fetched", "changes", changes)

	return NewChangelog(c.fromRef, c.toRef, c.repoURL, changes)
}

func (c Generator) changesFromCommits(commits []git.Commit) (Changes, error) {
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
	return issuesByEpic, nil
}

func (c Generator) fetchJiraIssue(commit git.Commit) (jira.Issue, error) {
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

func NewGenerator(jiraConfig jira.Config, fromRef, toRef, repoURL string) *Generator {
	client := jira.NewClient(jiraConfig)
	return &Generator{
		jiraConfig,
		fromRef,
		toRef,
		repoURL,
		client,
	}
}
