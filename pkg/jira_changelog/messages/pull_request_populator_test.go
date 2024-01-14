package messages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoPath_WhenRepoURLIsValid(t *testing.T) {
	repoURL := "https://github.com/HandOfGod94/jira_changelog"
	result, err := repoPath(repoURL)
	assert.NoError(t, err)
	assert.Equal(t, []string{"HandOfGod94", "jira_changelog"}, result)
}

func TestRepoPath_WhenRepoURLIsInvalid(t *testing.T) {
	repoURL := "https://google.com/foo?bar=fizz"
	result, err := repoPath(repoURL)
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo"}, result)
}
