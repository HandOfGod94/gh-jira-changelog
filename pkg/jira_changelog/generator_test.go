package jira_changelog

import (
	"testing"
	"time"

	"github.com/handofgod94/gh-jira-changelog/mocks"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/git"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/stretchr/testify/assert"
)

func TestChangelogFromCommits(t *testing.T) {
	commits := []git.Commit{
		{Time: time.Now(), Message: "[TEST-1234] commit message1", Sha: "3245vw"},
		{Time: time.Now(), Message: "[TEST-4546] commit message sample1", Sha: "3245vw"},
		{Time: time.Now(), Message: "[TEST-1234] commit message2", Sha: "3245vw"},
		{Time: time.Now(), Message: "[TEST-4546] commit message sample2", Sha: "3245vw"},
		{Time: time.Now(), Message: "[TEST-12345] commit message from same epic", Sha: "3245vw"},
		{Time: time.Now(), Message: "[NO-CARD] commit message random", Sha: "3245vw"},
		{Time: time.Now(), Message: "foobar commit message random", Sha: "3245vw"},
	}

	expected := &Changelog{
		Changes: map[string][]jira.Issue{
			"Epic1": {
				jira.NewIssue("TEST-1234", "Ticket description", "done", "Epic1"),
				jira.NewIssue("TEST-12345", "Ticket description of another from same epic", "done", "Epic1"),
			},
			"Epic2": {
				jira.NewIssue("TEST-4546", "Ticket description for 4546 issue", "done", "Epic2"),
			},
			"Misc": {
				jira.NewIssue("NO-CARD", "Ticket description for no card issue", "done", ""),
				jira.NewIssue("", "foobar commit message random", "", ""),
			},
		},
	}

	mockedClient := mocks.NewClient(t)
	mockedClient.On("FetchIssue", "TEST-1234").Return(jira.NewIssue("TEST-1234", "Ticket description", "done", "Epic1"), nil).Times(2)
	mockedClient.On("FetchIssue", "TEST-4546").Return(jira.NewIssue("TEST-4546", "Ticket description", "done", "Epic2"), nil).Times(2)
	mockedClient.On("FetchIssue", "TEST-12345").Return(jira.NewIssue("TEST-12345", "Ticket description", "done", "Epic1"), nil)
	mockedClient.On("FetchIssue", "NO-CARD").Return(jira.NewIssue("", "", "", ""), nil)
	generator := Generator{JiraConfig: jira.Config{ProjectName: "TEST"}}
	generator.client = mockedClient

	result := generator.changelogFromCommits(commits)

	assert.Equal(t, expected, result)
}
