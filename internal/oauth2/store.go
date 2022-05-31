package oauth2

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getToken() (tokenResponse, error) {
	dir := getNotionDir()
	token, err := os.ReadFile(filepath.Join(dir, "token"))
	if err != nil {
		return tokenResponse{}, err
	}
	res := tokenResponse{}
	err = json.Unmarshal(token, &res)
	return res, err
}

func storeToken(token string) error {
	dir := getNotionDir()
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(dir, "token"), []byte(token), 0644)
	return err
}

func getNotionDir() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".notion")
	return dir
}
