package messages

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

var gitoutputPattern = regexp.MustCompile(`^\((\d+)\)\s+\{(\w+)\}\s*(.*)`)

type GitOutput string

func (gt GitOutput) Commits() ([]Commit, error) {
	output := strings.TrimSpace(string(gt))
	commits := make([]Commit, 0)
	for _, line := range strings.Split(output, "\n") {
		message, err := extractCommitMessage(line)
		if err != nil {
			slog.Error("failed to extract commit message", "gitlogLine", line)
			return []Commit{}, fmt.Errorf("failed to extract commit message. %w", err)
		}

		commitTime, err := extractTime(line)
		if err != nil {
			slog.Error("failed to extract timestamp", "gitlogLine", line)
			return []Commit{}, fmt.Errorf("failed to extract timestamp. %w", err)
		}

		sha, err := extractSha(line)
		if err != nil {
			slog.Error("failed to extract sha", "gitlogLine", line)
			return []Commit{}, fmt.Errorf("failed to extract sha. %w", err)
		}

		commits = append(commits, Commit{
			Summary: message,
			Time:    commitTime,
			Sha:     sha,
		})
	}
	return commits, nil
}

func extractCommitMessage(gitlogLine string) (string, error) {
	gitlogLine = strings.TrimSpace(gitlogLine)
	result := gitoutputPattern.FindStringSubmatch(gitlogLine)
	if len(result) < 4 {
		return "", fmt.Errorf("couldn't find commit message in git log. %v", gitlogLine)
	}

	return result[3], nil
}

func extractTime(gitlogLine string) (time.Time, error) {
	result := gitoutputPattern.FindStringSubmatch(gitlogLine)
	if len(result) < 2 {
		return time.Time{}, fmt.Errorf("couldn't find timestamp in commit message. %v", gitlogLine)
	}

	timestamp, err := strconv.Atoi(result[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to extract timestamp in commit message. %w", err)
	}
	return time.Unix(int64(timestamp), 0), nil
}

func extractSha(gitlogLine string) (string, error) {
	result := gitoutputPattern.FindStringSubmatch(gitlogLine)
	if len(result) < 3 {
		return "", fmt.Errorf("couldn't find sha in commit message. %v", gitlogLine)
	}

	return result[2], nil
}
