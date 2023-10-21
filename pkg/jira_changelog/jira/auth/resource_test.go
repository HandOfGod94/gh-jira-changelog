package auth

import (
	"os"
	"path"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestDefaultDir_When_XDG_CONFIG_HOME_isSet(t *testing.T) {
	os.Setenv("XDG_CONFIG_HOME", "~/foobar")
	got, err := DefaultConfDir()

	assert.NoError(t, err)
	assert.Equal(t, "~/foobar/gh-jira-changelog", got)
}

func TestDefaultDir_When_XDG_CONFIG_HOME_isNotSet(t *testing.T) {
	os.Unsetenv("XDG_CONFIG_HOME")

	homeDir, _ := homedir.Dir()
	expected := path.Join(homeDir, "gh-jira-changelog")

	got, err := DefaultConfDir()

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParseResources(t *testing.T) {
	raw := `[
	  {
		"id": "1324a887-45db-1bf4-1e99-ef0ff456d421",
		"name": "Site name",
		"url": "https://your-domain.atlassian.net",
		"scopes": [
		  "write:jira-work",
		  "read:jira-user",
		  "manage:jira-configuration"
		],
		"avatarUrl": "https:\/\/site-admin-avatar-cdn.prod.public.atl-paas.net\/avatars\/240\/flag.png"
	  }
	]`

	got, err := parseResources([]byte(raw))

	want := []Resource{
		{
			CloudId: "1324a887-45db-1bf4-1e99-ef0ff456d421",
			Name:    "Site name",
			BaseURL: "https://your-domain.atlassian.net",
			Scopes: []string{
				"write:jira-work",
				"read:jira-user",
				"manage:jira-configuration",
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
