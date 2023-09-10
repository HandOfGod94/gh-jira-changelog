package jira_changelog

import (
	"context"
	"testing"
	"time"

	"github.com/handofgod94/gh-jira-changelog/mocks"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/git"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/stretchr/testify/assert"
)

func TestFetchJiraIssuesEvent(t *testing.T) {
	commits := []git.Commit{
		{Time: time.Now(), Message: "[TEST-1234] commit message1", Sha: "3245vw"},
		{Time: time.Now(), Message: "[TEST-4546] commit message sample1", Sha: "3245vw"},
		{Time: time.Now(), Message: "[TEST-1234] commit message2", Sha: "3245vw"},
		{Time: time.Now(), Message: "[TEST-4546] commit message sample2", Sha: "3245vw"},
		{Time: time.Now(), Message: "[TEST-12345] commit message from same epic", Sha: "3245vw"},
		{Time: time.Now(), Message: "[NO-CARD] commit message random", Sha: "3245vw"},
		{Time: time.Now(), Message: "foobar commit message random", Sha: "3245vw"},
	}

	jiraIssues := []jira.Issue{
		jira.NewIssue("TEST-1234", "Ticket description", "done", "Epic1"),
		jira.NewIssue("TEST-4546", "Ticket description for 4546 issue", "done", "Epic2"),
		jira.NewIssue("TEST-12345", "Ticket description of another card from same epic", "done", "Epic1"),
		jira.NewIssue("", "[NO-CARD] commit message random (3245vw)", "done", ""),
		jira.NewIssue("", "foobar commit message random (3245vw)", "done", ""),
	}

	mockedClient := mocks.NewClient(t)
	mockedClient.On("FetchIssue", "TEST-1234").Return(jiraIssues[0], nil).Twice()
	mockedClient.On("FetchIssue", "TEST-4546").Return(jiraIssues[1], nil).Twice()
	mockedClient.On("FetchIssue", "TEST-12345").Return(jiraIssues[2], nil)

	// Setup
	generator := NewGenerator(jira.Config{}, "fromRef", "toRef", "http://example-repo.com")
	generator.client = mockedClient
	generator.commits = commits
	generator.FSM.SetState(CommitsFetched)

	// invoke event
	err := generator.FSM.Event(context.Background(), FetchJiraIssues)

	assert.NoError(t, err)
	assert.Equal(t, len(jiraIssues), len(generator.jiraIssues))
	assert.Equal(t, jiraIssues, generator.jiraIssues)
	assert.Equal(t, generator.FSM.Current(), JiraIssuesFetched)
}

func TestRecordChangeLogEvent(t *testing.T) {
	issues := []jira.Issue{
		jira.NewIssue("TEST-1234", "Ticket description", "done", "Epic1"),
		jira.NewIssue("TEST-12345", "Ticket description of another from same epic", "done", "Epic1"),
		jira.NewIssue("TEST-4546", "Ticket description for 4546 issue", "done", "Epic2"),
		jira.NewIssue("", "[NO-CARD] commit message random (3245vw)", "done", ""),
		jira.NewIssue("", "foobar commit message random (3245vw)", "done", ""),
	}
	expected := Changes{
		"Epic1": {
			jira.NewIssue("TEST-1234", "Ticket description", "done", "Epic1"),
			jira.NewIssue("TEST-12345", "Ticket description of another from same epic", "done", "Epic1"),
		},
		"Epic2": {
			jira.NewIssue("TEST-4546", "Ticket description for 4546 issue", "done", "Epic2"),
		},
		"Miscellaneous": {
			jira.NewIssue("", "[NO-CARD] commit message random (3245vw)", "done", ""),
			jira.NewIssue("", "foobar commit message random (3245vw)", "done", ""),
		},
	}

	generator := NewGenerator(jira.Config{}, "fromRef", "toRef", "http://example-repo.com")
	generator.jiraIssues = issues
	generator.FSM.SetState(JiraIssuesFetched)
	generator.FSM.Event(context.Background(), RecordChanges)

	actual := generator.changes
	assert.Equal(t, expected, actual)
	assert.Equal(t, generator.FSM.Current(), ChangesRecorded)
}
