package git

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/looplab/fsm"
)

type Commit struct {
	Message string
	Time    time.Time
	Sha     string
}

var gitoutputPattern = regexp.MustCompile(`^\((\d+)\)\s+\{(\w+)\}\s*(.*)`)

const (
	InitialState    = "initial"
	CommandExecuted = "command_executed"
	OutuptParsed    = "output_parsed"
)

type commitParseWorkflow struct {
	fromRef   string
	toRef     string
	gitOutput GitOutput
	commits   []Commit
	FSM       *fsm.FSM
}

func NewCommitParseWorkflow(fromRef, toRef string) *commitParseWorkflow {
	cpw := &commitParseWorkflow{
		fromRef: fromRef,
		toRef:   toRef,
	}

	cpw.FSM = fsm.NewFSM(
		InitialState,
		fsm.Events{
			{Name: "execute_git_log", Src: []string{InitialState}, Dst: CommandExecuted},
			{Name: "parse_output", Src: []string{CommandExecuted}, Dst: OutuptParsed},
		},
		fsm.Callbacks{
			"before_execute_git_log": func(ctx context.Context, e *fsm.Event) {
				ouptut, err := execGitLog(ctx, fromRef, toRef)
				if err != nil {
					e.Cancel(err)
					return
				}
				cpw.gitOutput = ouptut
			},
			"before_parse_output": func(ctx context.Context, e *fsm.Event) {
				commits, err := cpw.gitOutput.Commits()
				if err != nil {
					e.Cancel(err)
					return
				}
				cpw.commits = commits
			},
		},
	)
	return cpw
}

func (cpw *commitParseWorkflow) Commits(ctx context.Context) ([]Commit, error) {
	err := cpw.FSM.Event(ctx, "execute_git_log")
	if err != nil {
		return []Commit{}, fmt.Errorf("failed to execute git log. %w", err)
	}

	err = cpw.FSM.Event(ctx, "parse_output")
	if err != nil {
		return []Commit{}, fmt.Errorf("failed to parse output. %w", err)
	}

	return cpw.commits, nil
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
