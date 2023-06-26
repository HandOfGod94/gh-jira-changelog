package jira

import (
	"regexp"

	"github.com/handofgod94/jira_changelog/pkg/jira_changelog/git"
)

type Config struct {
	BaseUrl     string
	ProjectName string
	User        string
	ApiToken    string
}

type Issue struct {
	Id     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		Parent struct {
			Fields struct {
				Summary string `json:"summary"`
			} `json:"fields"`
		} `json:"parent"`
		Status struct {
			StatusCategory struct {
				Key string `json:"key"`
			} `json:"statusCategory"`
		} `json:"status"`
		Summary string `json:"summary"`
	} `json:"fields"`
}

func NewIssue(key, summary string, status string) Issue {
	issues := &Issue{}
	issues.Key = key
	issues.Fields.Summary = summary
	issues.Fields.Status.StatusCategory.Key = status
	return *issues
}

func (i Issue) IsWip() bool {
	return i.Fields.Status.StatusCategory.Key != "done"
}

func (i Issue) Epic() string {
	if i.Fields.Parent.Fields.Summary != "" {
		return i.Fields.Parent.Fields.Summary
	}
	return ""
}

type JiraIssueId string

func IssueId(projectName string, commitMessage git.CommitMessage) JiraIssueId {
	jiraIssuePattern := regexp.MustCompile("(\\[)?" + projectName + "-(\\d*)(\\])?.*")
	result := jiraIssuePattern.FindStringSubmatch(string(commitMessage))
	if len(result) == 0 {
		return ""
	}
	return JiraIssueId(projectName + "-" + result[2])

}
