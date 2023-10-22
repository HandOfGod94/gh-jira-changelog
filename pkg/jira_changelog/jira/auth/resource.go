package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

const ResourcesFile = "resources.json"

type Resource struct {
	CloudId string   `json:"id"`
	Name    string   `json:"name"`
	BaseURL string   `json:"url"`
	Scopes  []string `json:"scopes"`
}

func (r Resource) Save() error {

	confdir, err := getOrCreateConfDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir %w", err)
	}

	filepath := path.Join(confdir, ResourcesFile)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(r)
	if err != nil {
		return fmt.Errorf("failed to encode resources to json. %w", err)
	}

	return nil
}

func parseResources(raw []byte) ([]Resource, error) {
	result := make([]Resource, 0)

	if err := json.Unmarshal(raw, &result); err != nil {
		return []Resource{}, nil
	}

	return result, nil
}

func DefaultConfDir() (res string, err error) {
	filepath := path.Join("gh-jira-changelog")
	res = os.Getenv("XDG_CONFIG_HOME")
	if res == "" {
		res, err = homedir.Dir()
		if err != nil {
			return
		}
	}

	return path.Join(res, filepath), nil
}

func getOrCreateConfDir() (string, error) {
	confdir, err := DefaultConfDir()
	if err != nil {
		return "", fmt.Errorf("failed to generate config file %w", err)
	}

	if _, err := os.Stat(confdir); os.IsNotExist(err) {
		err = os.Mkdir(confdir, os.ModeDir|0755)
		if err != nil {
			return "", err
		}
	}
	return confdir, nil
}
