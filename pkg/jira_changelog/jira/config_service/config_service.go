package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/exp/slog"
)

func Save(v any, filepath string) error {
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
	err = enc.Encode(v)
	if err != nil {
		return fmt.Errorf("failed to encode resources to json. %w", err)
	}

	return nil
}

func Load(v any, filepath string) (err error) {
	confdir, err := defaultConfDir()
	if err != nil {
		return
	}

	filepath = path.Join(confdir, filepath)
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(v)
	if err != nil {
		return
	}

	return nil
}

func Clear() error {
	confdir, err := defaultConfDir()
	if err != nil {
		return err
	}

	slog.Info("clearing config dir", "Dir", confdir)

	if err := os.RemoveAll(confdir); err != nil {
		slog.Error("failed to delete conf dir", "error", err)
		return err
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
