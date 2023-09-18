package git

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"time"
)

type Commit struct {
	Message string
	Time    time.Time
	Sha     string
}

var gitoutputPattern = regexp.MustCompile(`^\((\d+)\)\s+\{(\w+)\}\s*(.*)`)

type commitPopulator struct {
	fromRef string
	toRef   string
}

func NewCommitPopulator(fromRef, toRef string) *commitPopulator {
	cpw := &commitPopulator{
		fromRef: fromRef,
		toRef:   toRef,
	}
	return cpw
}

func (cpw *commitPopulator) Commits(ctx context.Context) ([]Commit, error) {
	gitOutput, err := execGitLog(ctx, cpw.fromRef, cpw.toRef)
	if err != nil {
		return []Commit{}, fmt.Errorf("failed to execute git log. %w", err)
	}

	commits, err := gitOutput.Commits()
	if err != nil {
		return []Commit{}, fmt.Errorf("failed to parse output. %w", err)
	}

	return commits, nil
}

func execGitLog(ctx context.Context, fromRef, toRef string) (GitOutput, error) {
	cmd := exec.CommandContext(ctx, "git", "log", "--decorate-refs-exclude=refs/*", "--pretty=(%ct) {%h} %d %s", "--no-merges", fromRef+".."+toRef)
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute git command: %v", err)
	}

	result := string(stdout)
	return GitOutput(result), nil
}
