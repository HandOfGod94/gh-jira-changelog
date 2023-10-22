package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

type Token struct {
	AccessToken  string    `json:"access_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
}

const TokenFile = "token.json"

func (t *Token) Save() error {
	confdir, err := getOrCreateConfDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir for saving token. %w", err)
	}

	filepath := path.Join(confdir, TokenFile)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(t)
	if err != nil {
		return fmt.Errorf("failed to encode resources to json. %w", err)
	}

	return nil
}
