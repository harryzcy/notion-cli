package oauth2

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var (
	notionDir = ""
)

func init() {
	home, _ := os.UserHomeDir()
	notionDir = filepath.Join(home, ".notion")
}

// GetToken returns the token previously stored
func GetToken() (tokenResponse, error) {
	token, err := os.ReadFile(filepath.Join(notionDir, "token"))
	if err != nil {
		return tokenResponse{}, err
	}
	res := tokenResponse{}
	err = json.Unmarshal(token, &res)
	return res, err
}

func storeToken(token string) error {
	err := os.MkdirAll(notionDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(notionDir, "token"), []byte(token), 0644)
	return err
}
