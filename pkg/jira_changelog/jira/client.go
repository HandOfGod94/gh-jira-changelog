package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/slog"
)

type Client interface {
	FetchIssue(issueId string) (Issue, error)
}

type client struct {
	config     Context
	httpClient *resty.Client
}

func (c *client) setupClient() {
	c.httpClient = resty.New()
	c.httpClient.SetBaseURL(c.config.BaseURL)
	c.httpClient.SetBasicAuth(c.config.User, c.config.ApiToken)
	c.httpClient.SetTimeout(5 * time.Second)
}

func (c *client) FetchIssue(issueId string) (Issue, error) {
	requestUrl, err := url.JoinPath(c.config.BaseURL, "rest", "api", "3", "issue", issueId)
	slog.Debug("Preparing fetch request", "url", requestUrl)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to create request url. %w", err)
	}

	resp, err := c.httpClient.R().Get(requestUrl)
	if err != nil || resp.StatusCode() != http.StatusOK {
		return Issue{}, fmt.Errorf("failed to fetch issue. code: %v, %w", resp.StatusCode(), err)
	}

	var issue Issue
	if err := json.Unmarshal(resp.Body(), &issue); err != nil {
		return Issue{}, fmt.Errorf("failed to decode issue. %w", err)
	}

	return issue, nil
}

func NewClient(config Context) Client {
	c := &client{config: config}
	c.setupClient()
	return c
}
