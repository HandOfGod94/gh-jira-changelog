package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	giturls "github.com/whilp/git-urls"
	"golang.org/x/exp/slog"
)

func CurrentRepoURL(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url")
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote origin url: %w", err)
	}

	result := strings.TrimSpace(string(stdout))
	url, err := giturls.Parse(result)
	if err != nil {
		slog.Error("failed to parse git url.", "error", err)
		return "", fmt.Errorf("failed to parse git url: %w", err)
	}

	host := url.Host
	user := strings.Split(url.Path, "/")[0]
	repoName := strings.TrimSuffix(strings.Split(url.Path, "/")[1], ".git")

	return fmt.Sprintf("https://%s/%s/%s", host, user, repoName), nil
}
