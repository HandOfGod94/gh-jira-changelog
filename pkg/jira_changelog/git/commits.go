package git

import (
	"fmt"
	"os/exec"
	"strings"
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

	rawCommitMessages := strings.Split(gitlogs, "\n")
	commitMessages := make([]CommitMessage, 0, len(rawCommitMessages))
	for _, rawCommitMessage := range rawCommitMessages {
		rawCommitMessage = strings.TrimSpace(rawCommitMessage)
		commitMessages = append(commitMessages, CommitMessage(rawCommitMessage))
		fmt.Println(rawCommitMessage)
	}
	return commitMessages, nil
}
