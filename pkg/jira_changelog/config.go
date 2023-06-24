package jira_changelog

type JiraConfig struct {
	ProjectUrl  string
	ProjectName string
	ApiToken    string
}

type GitConfig struct {
	FromRef string
	ToRef   string
}
