package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

func Save(content any, filepath string) error {
	confdir, err := getOrCreateConfDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir for saving token. %w", err)
	}

	filepath = path.Join(confdir, filepath)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(content)
	if err != nil {
		return fmt.Errorf("failed to encode resources to json. %w", err)
	}

	return nil
}

func defaultConfDir() (res string, err error) {
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
	confdir, err := defaultConfDir()
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
