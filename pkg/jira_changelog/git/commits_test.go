package git_test

import (
	"testing"
	"time"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/git"
	"github.com/stretchr/testify/assert"
)

func TestCommits(t *testing.T) {
	testCases := []struct {
		desc      string
		gitOutput git.GitOutput
		want      []git.Commit
		wantErr   bool
	}{
		{
			desc: "returns commits when gitoutput is valid",
			gitOutput: git.GitOutput(`
(1687839814)  use extra space while generating template
(1688059937)  [JIRA-123] refactor: extract out structs from jira/types
(1687799347)  add warning emoji for changelog lineitem
			`),
			want: []git.Commit{
				{
					Message: "use extra space while generating template",
					Time:    time.Unix(1687839814, 0),
				},
				{
					Message: "[JIRA-123] refactor: extract out structs from jira/types",
					Time:    time.Unix(1688059937, 0),
				},
				{
					Message: "add warning emoji for changelog lineitem",
					Time:    time.Unix(1687799347, 0),
				},
			},
		},
		{
			desc: "returns single commit if gitoutput has single line",
			gitOutput: git.GitOutput(`
(1688059937)  refactor: extract out structs from jira/types
`),
			want: []git.Commit{{Message: "refactor: extract out structs from jira/types", Time: time.Unix(1688059937, 0)}},
		},
		{
			desc:      "returns error when output is not in correct format",
			gitOutput: git.GitOutput(`foobar`),
			wantErr:   true,
		},
		{
			desc:      "returns error when output is empty",
			gitOutput: git.GitOutput(""),
			wantErr:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := tc.gitOutput.Commits()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
