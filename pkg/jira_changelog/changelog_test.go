package jira_changelog_test

import (
	"bytes"
	"testing"

	"github.com/handofgod94/jira_changelog/pkg/jira_changelog"
	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/jira"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	testCases := []struct {
		desc      string
		changelog jira_changelog.Changelog
		want      string
	}{
		{
			desc: "when there are `done` issues",
			changelog: jira_changelog.Changelog{
				Changes: map[string][]jira.Issue{
					"TestEpic": {
						jira.NewIssue("TEST-1", "foobar is new", "done"),
						jira.NewIssue("TEST-2", "fizzbuzz is something else", "done"),
					},
				},
			},
			want: `## What has changed?

### TestEpic
- [TEST-1] foobar is new
- [TEST-2] fizzbuzz is something else
`,
		},
		{
			desc: "when there are `wip` issues",
			changelog: jira_changelog.Changelog{
				Changes: map[string][]jira.Issue{
					"TestEpic": {
						jira.NewIssue("TEST-1", "foobar is new", "done"),
						jira.NewIssue("TEST-2", "fizzbuzz is something else", "in progress"),
						jira.NewIssue("TEST-3", "fizzbuzz is something else", "done"),
					},
				},
			},
			want: `## What has changed?

### TestEpic
- [TEST-1] foobar is new
- :warning: [TEST-2] fizzbuzz is something else
- [TEST-3] fizzbuzz is something else
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := bytes.NewBufferString("")
			tc.changelog.Render(got)
			assert.Equal(t, tc.want, got.String())
		})
	}
}
