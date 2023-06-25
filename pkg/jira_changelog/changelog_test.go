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
				DoneChanges: map[string][]jira.Issue{
					"TestEpic": {
						jira.NewIssue("TEST-1", "foobar is new"),
						jira.NewIssue("TEST-2", "fizzbuzz is something else"),
					},
				},
			},
			want: `## What has changed?

### TestEpic
- [TEST-1] foobar is new
- [TEST-2] fizzbuzz is something else

## :warning: WIP
### These cards are still in "done" state. Be careful while examining it
`,
		},
		// {
		// 	desc: "when there are `wip` issues",
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := bytes.NewBufferString("")
			tc.changelog.Render(got)
			assert.Equal(t, tc.want, got.String())
		})
	}
}
