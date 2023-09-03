package jira_changelog_test

import (
	"bytes"
	"testing"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
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
						jira.NewIssue("TEST-1", "foobar is new", "done", "TestEpic"),
						jira.NewIssue("TEST-2", "fizzbuzz is something else", "done", "TestEpic"),
					},
				},
			},
			want: `## What's changed?

### TestEpic
- [TEST-1] foobar is new
- [TEST-2] fizzbuzz is something else

:warning: = Work in Progress. Ensure that these cards don't break things in production.
`,
		},
		{
			desc: "when there are `wip` issues",
			changelog: jira_changelog.Changelog{
				Changes: map[string][]jira.Issue{
					"TestEpic": {
						jira.NewIssue("TEST-1", "foobar is new", "done", "TestEpic"),
						jira.NewIssue("TEST-2", "fizzbuzz is something else", "in progress", "TestEpic"),
						jira.NewIssue("TEST-3", "fizzbuzz is something else", "done", "TestEpic"),
					},
				},
			},
			want: `## What's changed?

### TestEpic
- [TEST-1] foobar is new
- :warning: [TEST-2] fizzbuzz is something else
- [TEST-3] fizzbuzz is something else

:warning: = Work in Progress. Ensure that these cards don't break things in production.
`,
		},
		{
			desc: "when there are multiple epics",
			changelog: jira_changelog.Changelog{
				Changes: map[string][]jira.Issue{
					"TestEpic1": {
						jira.NewIssue("TEST-1", "foobar is new", "done", "TestEpic1"),
						jira.NewIssue("TEST-2", "fizzbuzz is something else", "in progress", "TestEpic1"),
						jira.NewIssue("TEST-3", "fizzbuzz is something else", "done", "TestEpic1"),
					},
					"TestEpic2": {
						jira.NewIssue("TEST-4", "foobar is new", "done", "TestEpic2"),
						jira.NewIssue("TEST-5", "fizzbuzz is something else", "done", "TestEpic2"),
						jira.NewIssue("TEST-6", "fizzbuzz is something else", "done", "TestEpic2"),
					},
				},
			},
			want: `## What's changed?

### TestEpic1
- [TEST-1] foobar is new
- :warning: [TEST-2] fizzbuzz is something else
- [TEST-3] fizzbuzz is something else

### TestEpic2
- [TEST-4] foobar is new
- [TEST-5] fizzbuzz is something else
- [TEST-6] fizzbuzz is something else

:warning: = Work in Progress. Ensure that these cards don't break things in production.
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
