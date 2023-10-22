package oauth

import (
	"encoding/json"
	"fmt"

	config "github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira/config_service"
)

const ResourcesFile = "resources.json"

type Resource struct {
	CloudId string   `json:"id"`
	Name    string   `json:"name"`
	BaseURL string   `json:"url"`
	Scopes  []string `json:"scopes"`
}

func (r Resource) Save() error {
	return config.Save(r, ResourcesFile)
}

func parseResources(raw []byte) ([]Resource, error) {
	result := make([]Resource, 0)

	if err := json.Unmarshal(raw, &result); err != nil {
		return []Resource{}, nil
	}

	return result, nil
}

func LoadResource() (*Resource, error) {
	res := &Resource{}
	err := config.Load(res, ResourcesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load resources: %w", err)
	}

	return res, nil
}
