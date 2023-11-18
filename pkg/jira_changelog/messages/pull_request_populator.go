package messages

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/samber/lo"
	giturls "github.com/whilp/git-urls"
	"golang.org/x/exp/slog"
)

var _ Populator = &pullRequestPopulator{}
var _ Messager = &PullRequest{}

var prRegexPattern = `\* (?P<title>.+) by @(?P<author>\S+) in (?P<url>\S+)`

type PullRequest struct {
	Title  string
	Author string
	URL    string
}

func (p PullRequest) Message() string {
	return p.Title
}

type pullRequestPopulator struct {
	fromRef   string
	toRef     string
	apiClient *api.RESTClient
	repoOwner string
	repoName  string
}

func NewPullRequestPopulator(fromRef, toRef, repoURL string) (Populator, error) {
	apiClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}

	repoOwner, err := repoOwner(repoURL)
	if err != nil {
		return nil, err
	}

	repoName, err := repoName(repoURL)
	if err != nil {
		return nil, err
	}

	return &pullRequestPopulator{
		fromRef,
		toRef,
		apiClient,
		repoOwner,
		repoName,
	}, nil
}

func (p *pullRequestPopulator) Populate(ctx context.Context) ([]Messager, error) {
	pullRequests, err := p.PullRequests(ctx)
	if err != nil {
		return []Messager{}, err
	}

	messages := lo.Map(pullRequests, func(pullRequest PullRequest, i int) Messager { return pullRequest })
	return messages, nil
}

func (p *pullRequestPopulator) PullRequests(ctx context.Context) ([]PullRequest, error) {
	response := struct {
		Name string
		Body string
	}{}

	requestBody, err := json.Marshal(map[string]string{
		"owner":             p.repoOwner,
		"repo":              p.repoName,
		"tag_name":          p.toRef,
		"target_commitish":  "main", // TODO: make this configurable
		"previous_tag_name": p.fromRef,
	})
	if err != nil {
		return []PullRequest{}, err
	}

	slog.Debug("fetching changelog from github")

	err = p.apiClient.Post(fmt.Sprintf("repos/%s/%s/releases/generate-notes", p.repoOwner, p.repoName), bytes.NewBuffer(requestBody), &response)
	if err != nil {
		return []PullRequest{}, err
	}

	pullRequests, err := parsePullRequestBody(response.Body)
	if err != nil {
		slog.Error("error parsing pull request body", "error", err, "body", response.Body)
		return []PullRequest{}, err
	}

	slog.Info("successfully fetched changelog from github")
	slog.Debug("here's the changelog provided by github", "changelog", response.Body)

	return pullRequests, nil
}

func parsePullRequestBody(body string) ([]PullRequest, error) {
	pullrequests := make([]PullRequest, 0)
	lines := strings.Split(body, "\n")
	lines = lo.Map(lines, func(line string, i int) string { return strings.TrimSpace(line) })
	lines = lo.Filter(lines, func(line string, i int) bool { return line != "" })
	lines = lo.Filter(lines, func(line string, i int) bool { return !strings.HasPrefix(line, "## What's Changed") })
	lines = lo.Filter(lines, func(line string, i int) bool { return !strings.HasPrefix(line, "**Full Changelog**") })

	for _, line := range lines {
		pullrequest, err := parsePullRequestMessage(line)
		if err != nil {
			slog.Error("error parsing pull request message", "error", err, "line", line)
			return []PullRequest{}, err
		}
		pullrequests = append(pullrequests, pullrequest)
	}

	return pullrequests, nil
}

func parsePullRequestMessage(line string) (PullRequest, error) {
	re := regexp.MustCompile(prRegexPattern)
	result := re.FindStringSubmatch(line)
	if len(result) < 3 {
		return PullRequest{}, fmt.Errorf("invalid pull request title: %s", line)
	}

	title := re.SubexpIndex("title")
	author := re.SubexpIndex("author")
	url := re.SubexpIndex("url")
	return PullRequest{
		Title:  result[title],
		Author: result[author],
		URL:    result[url],
	}, nil
}

func repoOwner(repoURL string) (string, error) {
	url, err := giturls.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("error parsing repo url: %w", err)
	}

	path := strings.Split(url.Path, "/")
	if len(path) < 2 {
		return "", fmt.Errorf("invalid repo url: %s", repoURL)
	}

	return path[len(path)-2], nil
}

func repoName(repoURL string) (string, error) {
	url, err := giturls.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("error parsing repo url: %w", err)
	}

	path := strings.Split(url.Path, "/")
	if len(path) < 2 {
		return "", fmt.Errorf("invalid repo url: %s", repoURL)
	}

	return path[len(path)-1], nil
}
