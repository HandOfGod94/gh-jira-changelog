package messages

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/samber/lo"
)

var _ Populator = &commitPopulator{}
var _ Message = &Commit{}

type Commit struct {
	Summary string
	Time    time.Time
	Sha     string
}

func (c Commit) Message() string {
	return c.Summary
}

type commitPopulator struct {
	fromRef string
	toRef   string
}

func NewCommitPopulator(fromRef, toRef string) Populator {
	cpw := &commitPopulator{
		fromRef: fromRef,
		toRef:   toRef,
	}
	return cpw
}

func (cpw *commitPopulator) Populate(ctx context.Context) ([]Message, error) {
	gitOutput, err := execGitLog(ctx, cpw.fromRef, cpw.toRef)
	if err != nil {
		return nil, fmt.Errorf("failed to execute git log. %w", err)
	}

	commits, err := gitOutput.Commits()
	if err != nil {
		return nil, fmt.Errorf("failed to parse output. %w", err)
	}

	messages := lo.Map(commits, func(commit Commit, i int) Message { return commit })
	return messages, nil
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
