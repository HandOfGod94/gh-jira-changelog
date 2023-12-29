package git

import (
	"context"
	"os/exec"
)

func ConfigRemoteOriginCommand(ctx context.Context) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url")
	return cmd.Output()
}

func GitLogCommand(ctx context.Context, fromRef, toRef string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", "log", "--decorate-refs-exclude=refs/*", "--pretty=(%ct) {%h} %d %s", "--no-merges", fromRef+".."+toRef)
	return cmd.Output()
}
