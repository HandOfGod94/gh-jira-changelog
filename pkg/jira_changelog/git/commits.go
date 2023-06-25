package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/samber/lo"
)

type CommitMessage string

func CommitMessages(fromRef, toRef string) ([]CommitMessage, error) {
	cmd := exec.Command("git", "log", "--decorate-refs-exclude=refs/tags", "--pretty='%d %s'", "--no-merges", fromRef+".."+toRef)
	stdout, err := cmd.Output()
	if err != nil {
		return []CommitMessage{}, fmt.Errorf("failed to execute git command: %v", err)
	}

	gitlogs := string(stdout)
	gitlogs = strings.TrimSpace(gitlogs)

	commitMessages := lo.Map(strings.Split(gitlogs, "\n"), func(commitMessage string, index int) CommitMessage {
		commitMessage = strings.TrimSpace(commitMessage)
		return CommitMessage(commitMessage)
	})
	return commitMessages, nil
}
