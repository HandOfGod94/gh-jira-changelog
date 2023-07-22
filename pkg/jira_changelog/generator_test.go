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
			"Miscellaneous": {
				jira.NewIssue("", "[NO-CARD] commit message random", "done", ""),
				jira.NewIssue("", "foobar commit message random", "done", ""),
			},
		},
	}

	mockedClient := mocks.NewClient(t)
	mockedClient.On("FetchIssue", "TEST-1234").Return(jira.NewIssue("TEST-1234", "Ticket description", "done", "Epic1"), nil).Twice()
	mockedClient.On("FetchIssue", "TEST-4546").Return(jira.NewIssue("TEST-4546", "Ticket description for 4546 issue", "done", "Epic2"), nil).Twice()
	mockedClient.On("FetchIssue", "TEST-12345").Return(jira.NewIssue("TEST-12345", "Ticket description of another from same epic", "done", "Epic1"), nil)
	generator := Generator{JiraConfig: jira.Config{ProjectName: "TEST"}}
	generator.client = mockedClient

	result, err := generator.changelogFromCommits(commits)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
