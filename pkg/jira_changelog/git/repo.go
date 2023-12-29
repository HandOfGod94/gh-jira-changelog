package git

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	giturls "github.com/whilp/git-urls"
	"golang.org/x/exp/slog"
)

func CurrentRepoURL(ctx context.Context) (string, error) {
	resultBytes, err := ConfigRemoteOriginCommand(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get remote origin url: %w", err)
	}

	u, err := gitURLtoHttpURL(string(resultBytes))
	if err != nil {
		return "", fmt.Errorf("failed to convert git URL to HTTP URL: %w", err)
	}

	return u, nil
}

func gitURLtoHttpURL(gitURL string) (string, error) {
	guRaw := strings.TrimSpace(gitURL)
	gu, err := giturls.Parse(guRaw)
	if err != nil {
		slog.Error("failed to parse git url.", "error", err)
		return "", fmt.Errorf("failed to parse git url: %w", err)
	}

	paths := strings.Split(gu.Path, "/")

	if len(paths) < 2 {
		return "", fmt.Errorf("invalid git url provided: %w", err)
	}

	host := gu.Host
	user := paths[0]
	repoName := strings.TrimSuffix(paths[1], ".git")

	repoPath, err := url.JoinPath(user, repoName)
	if err != nil {
		return "", err
	}

	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   repoPath,
	}

	return u.String(), nil
}
