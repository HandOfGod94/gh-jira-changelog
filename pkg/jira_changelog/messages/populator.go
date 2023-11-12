package messages

import (
	"context"

	"golang.org/x/exp/slog"
)

type Message interface {
	Message() string
}

type Populator interface {
	Populate(ctx context.Context) ([]Message, error)
}

func NewCommitOrPRPopualtor(usePR bool, fromRef, toRef, repoURL string) (Populator, error) {
	if usePR {
		slog.Debug("using github PR titles to generate changelog")
		return NewPullRequestPopulator(fromRef, toRef, repoURL)
	} else {
		slog.Debug("using commit messages to generate changelog")
		return NewCommitPopulator(fromRef, toRef)
	}
}
