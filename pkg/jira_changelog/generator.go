package jira_changelog

import (
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/git"
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/jira"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
)

type Generator struct {
	JiraConfig jira.Config
	fromRef    string
	toRef      string
	client     *jira.Client
}

func (c Generator) Generate() {
	commitMessages, err := git.CommitMessages(c.fromRef, c.toRef)
	if err != nil {
		panic(err)
	}
	slog.Debug("Total commit messages", "count", len(commitMessages))

	jiraIssueIds := lo.Map(commitMessages, func(commitMessage git.CommitMessage, index int) jira.JiraIssueId {
		return jira.FromCommitMessage(c.JiraConfig.ProjectName, commitMessage)
	})
	jiraIssueIds = lo.Filter(jiraIssueIds, func(jiraIssueId jira.JiraIssueId, index int) bool { return jiraIssueId != "" })
	jiraIssueIds = lo.Uniq(jiraIssueIds)
	slog.Debug("Total jira issues ids", "count", len(jiraIssueIds))

	issues := lo.Map(jiraIssueIds, func(jiraIssueId jira.JiraIssueId, index int) jira.Issue {
		issue, err := c.client.FetchIssue(string(jiraIssueId))
		if err != nil {
			slog.Error("Error fetching issue", "issue", jiraIssueId, "error", err)
			panic(err)
		}
		slog.Debug("Fetched issue", "issue", issue)
		return *issue
	})
	slog.Debug("Total issues", "count", len(issues))

	issuesByEpic := lo.GroupBy(issues, func(issue jira.Issue) string { return issue.Epic() })
	slog.Info("Issues grouped by epic", "issues", issuesByEpic)
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
