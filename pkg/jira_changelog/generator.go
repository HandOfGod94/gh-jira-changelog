package jira_changelog

import "github.com/handofgod94/jira_changelog/pkg/jira_changelog/jira"

type Changelog struct {
	JiraConfig jira.Config
	FromRef    string
	ToRef      string
}

func (c Changelog) Generate() {

}
