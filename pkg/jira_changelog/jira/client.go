package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/exp/slog"
)

type Client struct {
	config     Config
	httpClient *http.Client
}

type Issue struct {
	Id     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		Parent struct {
			Fields struct {
				Summary string `json:"summary"`
			} `json:"fields,omitempty"`
		} `json:"parent,omitempty"`
	} `json:"fields"`
}

func (c *Client) setupClient() {
	c.httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
}

func (c *Client) attachDefaultHeaders(r *http.Request) {
	r.Header.Add("Accept", "application/json")
	r.SetBasicAuth(c.config.User, c.config.ApiToken)
}

func (c *Client) FetchIssue(issueId string) (*Issue, error) {
	requestUrl, err := url.JoinPath(c.config.BaseUrl, "rest", "api", "3", "issue", issueId)
	slog.Debug("Preparing fetch request", "url", requestUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create request url. %w", err)
	}

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request. %w", err)
	}
	c.attachDefaultHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue. %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch issue. status code: %d", resp.StatusCode)
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode issue. %w", err)
	}

	return &issue, nil
}

func NewClient(config Config) *Client {
	c := &Client{config: config}
	c.setupClient()
	return c
}
