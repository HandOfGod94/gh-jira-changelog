package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/slog"
)

type Client interface {
	FetchIssue(issueId string) (Issue, error)
}

type client struct {
	clientOpts *ClientOptions
	httpClient *resty.Client
}

func (c *client) FetchIssue(issueId string) (Issue, error) {
	requestUrl, err := url.JoinPath(c.clientOpts.BaseURL(), "rest", "api", "3", "issue", issueId)
	slog.Debug("Preparing fetch request", "url", requestUrl)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to create request url. %w", err)
	}

	resp, err := c.httpClient.R().Get(requestUrl)
	if err != nil || resp.StatusCode() != http.StatusOK {
		slog.Warn("failed to fetch issue", "code", resp.StatusCode(), "error", err)
		return Issue{}, fmt.Errorf("failed to fetch issue. code: %v, %w", resp.StatusCode(), err)
	}

	var issue Issue
	if err := json.Unmarshal(resp.Body(), &issue); err != nil {
		return Issue{}, fmt.Errorf("failed to decode issue. %w", err)
	}

	return issue, nil
}

func NewClient(clientOpts *ClientOptions) Client {
	c := &client{clientOpts: clientOpts}
	c.httpClient = clientOpts.Client()
	return c
}
