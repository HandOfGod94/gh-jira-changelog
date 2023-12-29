package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitURLToHttpURL_WhenGitURLIsValid(t *testing.T) {
	gitURL := "git@github.com:HandOfGod94/jira_changelog.git"

	result, err := gitURLtoHttpURL(gitURL)
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/HandOfGod94/jira_changelog", result)
}

func TestGitURLToHttpURL_WhenGitURLIsInvalid(t *testing.T) {
	gitURL := "foobar"

	_, err := gitURLtoHttpURL(gitURL)
	assert.Error(t, err)
}
