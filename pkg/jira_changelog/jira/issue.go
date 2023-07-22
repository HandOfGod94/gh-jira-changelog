package jira

import "regexp"

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

func NewIssue(key, summary, status, epic string) Issue {
	issue := &Issue{}
	issue.Key = key
	issue.Fields.Summary = summary
	issue.Fields.Status.StatusCategory.Key = status
	issue.Fields.Parent.Fields.Summary = epic
	return *issue
}

func (i Issue) IsWip() bool {
	return i.Fields.Status.StatusCategory.Key != "done"
}

func (i Issue) Epic() string {
	if i.Fields.Parent.Fields.Summary != "" {
		return i.Fields.Parent.Fields.Summary
	}
	return "Miscellaneous"
}

type JiraIssueId string

func IssueId(projectName, text string) JiraIssueId {
	jiraIssuePattern := regexp.MustCompile("(\\[)?" + projectName + "-(\\d*)(\\])?.*")
	result := jiraIssuePattern.FindStringSubmatch(text)
	if len(result) == 0 {
		return ""
	}
	return JiraIssueId(projectName + "-" + result[2])

}
