package jira_test

import (
	"testing"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"github.com/stretchr/testify/assert"
)

func TestIssueId(t *testing.T) {
	testCases := []struct {
		desc          string
		commitMessage string
		want          jira.JiraIssueId
	}{
		{
			desc:          "when jira issue is present in commit message",
			commitMessage: "[TEST-123] Test commit message",
			want:          jira.JiraIssueId("TEST-123"),
		},
		{
			desc:          "when jira issue is not present in commit message",
			commitMessage: "Test commit message",
			want:          jira.JiraIssueId(""),
		},
		{
			desc:          "when jira issue is of different project",
			commitMessage: "[OTHER-1235]Test commit message",
			want:          jira.JiraIssueId("OTHER-1235"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := jira.IssueId(tc.commitMessage)
			assert.Equal(t, tc.want, got)
		})
	}
}
