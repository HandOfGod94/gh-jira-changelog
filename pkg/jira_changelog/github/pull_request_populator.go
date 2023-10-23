package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"golang.org/x/exp/slog"
)

type PullRequest struct {
	Title  string
	Author string
	Number int
}

type pullRequestPopulator struct {
	fromRef   string
	toRef     string
	apiClient *api.RESTClient
	repoOwner string
	repoName  string
}

func NewPullRequestPopulator(fromRef, toRef, repoOwner, repoName string) (*pullRequestPopulator, error) {
	apiClient, err := api.DefaultRESTClient()
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

	err = p.apiClient.Post(fmt.Sprintf("repos/%s/%s/releases/generate-notes", p.repoOwner, p.repoName), bytes.NewBuffer(requestBody), &response)
	if err != nil {
		return []PullRequest{}, err
	}

	slog.Info("successfully fetched changelog from github")
	slog.Debug("here's the changelog provided by github", "changelog", response.Body)

	return []PullRequest{}, nil
}
