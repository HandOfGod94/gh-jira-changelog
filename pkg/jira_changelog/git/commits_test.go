package git_test

import (
	"testing"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/fsm_util"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/git"
	"github.com/stretchr/testify/assert"
)

func TestCommitParseWorkflow_Events(t *testing.T) {
	testCases := []struct {
		desc          string
		currentState  string
		allowedEvents []string
	}{
		{
			desc:          "stateTransition: InitialState -> CommandExecuted",
			currentState:  git.InitialState,
			allowedEvents: []fsm_util.Event{git.ExecuteGitLog},
		},
		{
			desc:          "stateTransition: CommandExecuted -> OutuptParsed",
			currentState:  git.CommandExecuted,
			allowedEvents: []fsm_util.Event{git.ParseOutput},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cpw := git.NewCommitParseWorkflow("fromRef", "toRef")
			cpw.FSM.SetState(tc.currentState)
			for _, transition := range tc.allowedEvents {
				assert.True(t, cpw.FSM.Can(transition))
			}
		})
	}
}
