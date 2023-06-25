package jira_test

import (
	"testing"

	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/git"
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/jira"
	"github.com/stretchr/testify/assert"
)

func TestIssueId(t *testing.T) {
	testCases := []struct {
		desc          string
		commitMessage git.CommitMessage
		want          jira.JiraIssueId
	}{
		{
			desc:          "when jira issue is present in commit message",
			commitMessage: git.CommitMessage("[TEST-123] Test commit message"),
			want:          jira.JiraIssueId("TEST-123"),
		},
		{
			desc:          "when jira issue is present but not in correct format",
			commitMessage: git.CommitMessage("TEST-123 Test commit message"),
			want:          jira.JiraIssueId("TEST-123"),
		},
		{
			desc:          "when jira issue is not present in commit message",
			commitMessage: git.CommitMessage("Test commit message"),
			want:          jira.JiraIssueId(""),
		},
		{
			desc:          "when jira issue is present in between in commit message",
			commitMessage: git.CommitMessage("[somethin-odd-1][TEST-1235]Test commit message"),
			want:          jira.JiraIssueId("TEST-1235"),
		},
		{
			desc:          "when jira issue is of different project",
			commitMessage: git.CommitMessage("[somethin-odd-1][OTHER-1235]Test commit message"),
			want:          jira.JiraIssueId(""),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := jira.IssueId("TEST", tc.commitMessage)
			assert.Equal(t, tc.want, got)
		})
	}
}
