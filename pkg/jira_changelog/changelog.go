package jira_changelog

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"text/template"

	"github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	"golang.org/x/exp/slog"
)

//go:embed templates/*.tmpl
var changeLogTmpl embed.FS

type Epic = string
type Changes = map[Epic][]jira.Issue

type Changelog struct {
	Changes Changes
	fromRef string
	toRef   string
	repoURL string
}

func (c *Changelog) Render(w io.Writer) {
	slog.Info("rendering changelog", "number-of-changes", len(c.Changes))
	tmpl, err := template.ParseFS(changeLogTmpl, "templates/changelog.tmpl")
	if err != nil {
		slog.Error("error parsing template", "error", err)
		panic(err)
	}

	resultBuffer := bytes.NewBufferString("")
	if err := tmpl.Execute(resultBuffer, c); err != nil {
		slog.Error("error executing template", "error", err)
		panic(err)
	}

	fmt.Fprint(w, resultBuffer.String())
}

func (c *Changelog) URL() string {
	return fmt.Sprintf("%s/compare/%s...%s", c.repoURL, c.fromRef, c.toRef)
}

func NewChangelog(fromRef, toRef, repoURL string, changes Changes) *Changelog {
	return &Changelog{
		Changes: changes,
		fromRef: fromRef,
		toRef:   toRef,
		repoURL: repoURL,
	}
}
