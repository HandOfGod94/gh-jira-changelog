package jira_changelog

type Changelog struct {
	JiraConfig JiraConfig
	GitConfig  GitConfig
}

func (c Changelog) Generate() {

}
