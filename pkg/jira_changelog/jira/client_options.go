package jira

import (
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira/oauth"
	"golang.org/x/exp/slog"
)

type ClientOptions struct {
	baseURL    string
	user       string
	apiToken   string
	oauthToken *oauth.Token
	resource   *oauth.Resource
	useOauth   bool
}

type Option string
type Options map[Option]string

const (
	BaseURL  Option = "base_url"
	User     Option = "user"
	ApiToken Option = "api_token"
)

func NewClientOptions(opts Options) *ClientOptions {
	c := &ClientOptions{
		baseURL:  opts[BaseURL],
		user:     opts[User],
		apiToken: opts[ApiToken],
	}

	token, err := oauth.LoadOauthToken()
	if err != nil {
		slog.Warn("failed to load oauth token. defaulting to api_token", "error", err)
	}
	c.oauthToken = token

	resource, err := oauth.LoadResource()
	if err != nil {
		slog.Warn("failed to load oauth resource. defaulting to api_token", "error", err)
	}
	c.resource = resource

	if c.apiToken == "" && c.oauthToken != nil && c.resource != nil {
		c.useOauth = true
	}

	return c
}

func (c *ClientOptions) BaseURL() string {
	if c.useOauth {
		res, err := url.JoinPath("https://api.atlassian.com", "ex", "jira", c.resource.CloudId)
		if err != nil {
			slog.Warn("failed to join oauth base url", "error", err)
			return ""
		}
		return res
	}
	return c.baseURL
}

func (c *ClientOptions) Client() *resty.Client {
	client := resty.New()

	if c.useOauth {
		client.SetAuthToken(c.oauthToken.AccessToken)
	} else {
		client.SetBasicAuth(c.user, c.apiToken)
	}

	client.SetBaseURL(c.BaseURL())
	client.SetTimeout(5 * time.Second)
	return client
}
