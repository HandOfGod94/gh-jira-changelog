package oauth

import (
	"fmt"
	"time"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira/config"
)

type Token struct {
	AccessToken  string    `json:"access_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
}

const TokenFile = "token.json"

func (t *Token) Save() error {
	return config.Save(t, TokenFile)
}

func LoadOauthToken() (*Token, error) {
	t := &Token{}
	err := config.Load(t, TokenFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load token file. %w", err)
	}
	return t, nil
}