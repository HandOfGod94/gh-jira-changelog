package messages

import (
	"context"
	"fmt"
	"time"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/git"
	"github.com/samber/lo"
)

var _ Populator = &commitPopulator{}
var _ Messager = &Commit{}

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

func NewCommitPopulator(fromRef, toRef string) (Populator, error) {
	cpw := &commitPopulator{
		fromRef: fromRef,
		toRef:   toRef,
	}
	return cpw, nil
}

func (cpw *commitPopulator) Populate(ctx context.Context) ([]Messager, error) {
	gitOutput, err := execGitLog(ctx, cpw.fromRef, cpw.toRef)
	if err != nil {
		return nil, fmt.Errorf("failed to execute git log. %w", err)
	}

	commits, err := gitOutput.Commits()
	if err != nil {
		return nil, fmt.Errorf("failed to parse output. %w", err)
	}

	messages := lo.Map(commits, func(commit Commit, i int) Messager { return commit })
	return messages, nil
}

func execGitLog(ctx context.Context, fromRef, toRef string) (GitOutput, error) {
	resultBytes, err := git.GitLogCommand(ctx, fromRef, toRef)
	if err != nil {
		return "", fmt.Errorf("failed to execute git command: %v", err)
	}

	result := string(resultBytes)
	return GitOutput(result), nil
}
