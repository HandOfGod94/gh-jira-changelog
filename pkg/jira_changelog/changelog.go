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

type Changelog struct {
	Changes map[string][]jira.Issue
}

func (c *Changelog) Render(w io.Writer) {
	slog.Info("rendering changelog")
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
