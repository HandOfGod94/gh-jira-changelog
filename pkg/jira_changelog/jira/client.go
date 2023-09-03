package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/exp/slog"
)

type Client interface {
	FetchIssue(issueId string) (Issue, error)
}

type client struct {
	config     Config
	httpClient *http.Client
}

func (c *client) setupClient() {
	c.httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
}

func (c *client) attachDefaultHeaders(r *http.Request) {
	r.Header.Add("Accept", "application/json")
	r.SetBasicAuth(c.config.User, c.config.ApiToken)
}

func (c *client) FetchIssue(issueId string) (Issue, error) {
	requestUrl, err := url.JoinPath(c.config.BaseURL, "rest", "api", "3", "issue", issueId)
	slog.Debug("Preparing fetch request", "url", requestUrl)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to create request url. %w", err)
	}

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to create request. %w", err)
	}
	c.attachDefaultHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to fetch issue. %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Issue{}, fmt.Errorf("failed to fetch issue. status code: %d", resp.StatusCode)
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return Issue{}, fmt.Errorf("failed to decode issue. %w", err)
	}

	return issue, nil
}

func NewClient(config Config) Client {
	c := &client{config: config}
	c.setupClient()
	return c
}
