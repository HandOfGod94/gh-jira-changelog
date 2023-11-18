package jira_changelog

import (
	"testing"
	"time"

	"github.com/handofgod94/gh-jira-changelog/mocks"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/messages"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestFetchJiraIssuesEvent(t *testing.T) {
	commits := []messages.Commit{
		{Time: time.Now(), Summary: "[TEST-1234] commit message1", Sha: "3245vw"},
		{Time: time.Now(), Summary: "[TEST-4546] commit message sample1", Sha: "3245vw"},
		{Time: time.Now(), Summary: "[TEST-1234] commit message2", Sha: "3245vw"},
		{Time: time.Now(), Summary: "[TEST-4546] commit message sample2", Sha: "3245vw"},
		{Time: time.Now(), Summary: "[TEST-12345] commit message from same epic", Sha: "3245vw"},
		{Time: time.Now(), Summary: "[NO-CARD] commit message random", Sha: "3245vw"},
		{Time: time.Now(), Summary: "foobar commit message random", Sha: "3245vw"},
	}

	want := []jira.Issue{
		jira.NewIssue("TEST-1234", "Ticket description", "done", "Epic1"),
		jira.NewIssue("TEST-4546", "Ticket description for 4546 issue", "done", "Epic2"),
		jira.NewIssue("TEST-12345", "Ticket description of another card from same epic", "done", "Epic1"),
		jira.NewIssue("", "[NO-CARD] commit message random", "done", ""),
		jira.NewIssue("", "foobar commit message random", "done", ""),
	}

	mockedClient := mocks.NewClient(t)
	mockedClient.On("FetchIssue", "TEST-1234").Return(want[0], nil).Twice()
	mockedClient.On("FetchIssue", "TEST-4546").Return(want[1], nil).Twice()
	mockedClient.On("FetchIssue", "TEST-12345").Return(want[2], nil)

	generator := NewGenerator(jira.NewClient(jira.NewClientOptions(nil)), false, "fromRef", "toRef", "http://example-repo.com")
	generator.client = mockedClient

	changeMessages := lo.Map(commits, func(commit messages.Commit, i int) messages.Messager { return commit })
	got, err := generator.fetchJiraIssues(changeMessages)

	assert.NoError(t, err)
	assert.Equal(t, len(want), len(got))
	assert.Equal(t, want, got)
}
