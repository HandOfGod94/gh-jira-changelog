package jira_changelog

import (
	"context"
	"fmt"

	. "github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/fsm_util"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/git"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/looplab/fsm"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
)

type Generator struct {
	JiraConfig jira.Config
	fromRef    string
	toRef      string
	repoURL    string
	client     jira.Client

	commits    []git.Commit
	jiraIssues []jira.Issue
	changes    Changes
	FSM        *fsm.FSM
}

const (
	Initial           = State("initial")
	CommitsFetched    = State("commits_fetched")
	JiraIssuesFetched = State("jira_issues_fetched")
	ChangelogRecored  = State("changelog_recorded")

	FetchCommits    = Event("fetch_commits")
	FetchJiraIssues = Event("fetch_jira_issues")
	RecordChangelog = Event("record_changelog")
)

func NewGenerator(jiraConfig jira.Config, fromRef, toRef, repoURL string) *Generator {
	client := jira.NewClient(jiraConfig)
	g := &Generator{
		JiraConfig: jiraConfig,
		fromRef:    fromRef,
		toRef:      toRef,
		repoURL:    repoURL,
		client:     client,
	}

	g.FSM = fsm.NewFSM(
		Initial,
		fsm.Events{
			{Name: FetchCommits, Src: []string{Initial}, Dst: CommitsFetched},
			{Name: FetchJiraIssues, Src: []string{CommitsFetched}, Dst: JiraIssuesFetched},
			{Name: RecordChangelog, Src: []string{JiraIssuesFetched}, Dst: ChangelogRecored},
		},
		fsm.Callbacks{
			Before(FetchCommits): func(ctx context.Context, e *fsm.Event) {
				gcw := git.NewCommitParseWorkflow(fromRef, toRef)
				commits, err := gcw.Commits(ctx)
				if err != nil {
					e.Cancel(err)
					return
				}
				g.commits = commits
			},
			Before(FetchJiraIssues): func(ctx context.Context, e *fsm.Event) {
				issues, err := g.fetchJiraIssues(g.commits)
				if err != nil {
					e.Cancel(err)
					return
				}
				g.jiraIssues = lo.Uniq(issues)
			},
			RecordChangelog: func(ctx context.Context, e *fsm.Event) {
				slog.Debug("Total jira issues ids", "count", len(g.jiraIssues))
				slog.Debug("Recroding changelog")

				issuesByEpic := lo.GroupBy(g.jiraIssues, func(issue jira.Issue) string { return issue.Epic() })

				slog.Debug("Total epics", "count", len(issuesByEpic))
				g.changes = issuesByEpic
			},
		},
	)

	return g
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *Generator) Generate(ctx context.Context) *Changelog {
	panicIfErr(c.FSM.Event(ctx, FetchCommits))
	panicIfErr(c.FSM.Event(ctx, FetchJiraIssues))
	panicIfErr(c.FSM.Event(ctx, RecordChangelog))

	slog.Debug("Total changelog items", "count", len(c.changes))

	return NewChangelog(c.fromRef, c.toRef, c.repoURL, c.changes)
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
	return jiraIssues, nil
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
