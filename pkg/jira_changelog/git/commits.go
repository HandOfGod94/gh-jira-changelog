package git

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	. "github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/fsm_util"
	"github.com/looplab/fsm"
)

type Commit struct {
	Message string
	Time    time.Time
	Sha     string
}

var gitoutputPattern = regexp.MustCompile(`^\((\d+)\)\s+\{(\w+)\}\s*(.*)`)

const (
	InitialState    = State("initial_state")
	CommandExecuted = State("command_executed")
	OutuptParsed    = State("output_parsed")

	ExecuteGitLog = Event("execute_git_log")
	ParseOutput   = Event("parse_output")
)

type commitParseWorkflow struct {
	fromRef   string
	toRef     string
	GitOutput GitOutput
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
			{Name: ExecuteGitLog, Src: []string{InitialState}, Dst: CommandExecuted},
			{Name: ParseOutput, Src: []string{CommandExecuted}, Dst: OutuptParsed},
		},
		fsm.Callbacks{
			Before(ExecuteGitLog): func(ctx context.Context, e *fsm.Event) {
				ouptut, err := execGitLog(ctx, fromRef, toRef)
				if err != nil {
					e.Cancel(err)
					return
				}
				cpw.GitOutput = ouptut
			},
			Before(ParseOutput): func(ctx context.Context, e *fsm.Event) {
				commits, err := cpw.GitOutput.Commits()
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
	err := cpw.FSM.Event(ctx, ExecuteGitLog)
	if err != nil {
		return []Commit{}, fmt.Errorf("failed to execute git log. %w", err)
	}

	err = cpw.FSM.Event(ctx, ParseOutput)
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
