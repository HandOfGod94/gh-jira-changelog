package config

import (
	"os"
	"path"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestDefaultDir_When_XDG_CONFIG_HOME_isSet(t *testing.T) {
	os.Setenv("XDG_CONFIG_HOME", "~/foobar")
	got, err := defaultConfDir()

	assert.NoError(t, err)
	assert.Equal(t, "~/foobar/gh-jira-changelog", got)
}

func TestDefaultDir_When_XDG_CONFIG_HOME_isNotSet(t *testing.T) {
	os.Unsetenv("XDG_CONFIG_HOME")

	homeDir, _ := homedir.Dir()
	expected := path.Join(homeDir, "gh-jira-changelog")

	got, err := defaultConfDir()

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
