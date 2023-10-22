package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
