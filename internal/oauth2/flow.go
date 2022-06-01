package oauth2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/browser"
)

var (
	baseURL = "https://api.notion.com/v1/oauth/authorize"
	port    = "30001"
)

func Flow(clientID, clientSecret string) error {
	toAuth := shouldAuth()
	if !toAuth {
		return nil
	}

	c := make(chan codeResponse)
	startErr := make(chan error)
	go func() {
		err := startServer(c)
		startErr <- err
	}()
	if err := <-startErr; err != nil {
		return fmt.Errorf("failed to start server")
	}

	state := generateRandomState()

	url := buildAuthURL(clientID, state)
	_ = browser.OpenURL(url)

	res := <-c
	valid := validateCodeResponse(res, state)
	if !valid {
		return fmt.Errorf("failed to obtain code grant")
	}

	err := exchangeGrant(clientID, clientSecret, res.code, state)
	if err != nil {
		return err
	}

	return nil
}

var scanln = fmt.Scanln

func shouldAuth() bool {
	token, err := GetToken()
	if err == nil && token.AccessToken != "" {
		// code is stored, ask if user wants to re-auth
		fmt.Print("Do you want to re-auth? (y/N): ")
		var input string
		_, err = scanln(&input)
		if err != nil {
			return false
		}
		return strings.ToLower(input) == "y"
	}
	return true
}

func validateCodeResponse(res codeResponse, state string) bool {
	if res.err != "" {
		fmt.Println("Access is not granted")
		return false
	}

	if res.state != state || res.state == "" {
		fmt.Println("State is not matched")
		return false
	}
	return true
}

func exchangeGrant(clientID, clientSecret, code, state string) error {
	url := "https://api.notion.com/v1/oauth/token"

	data := map[string]string{
		"grant_type":   "authorization_code",
		"code":         code,
		"redirect_uri": "http://localhost:" + port + "/callback",
	}
	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to exchange code grant")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", basicAuth(clientID, clientSecret))

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to exchange code grant")
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	var token tokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		return fmt.Errorf("failed to exchange code grant")
	}

	err = storeToken(string(body))
	if err != nil {
		return fmt.Errorf("failed to store token")
	}

	return nil
}

func basicAuth(clientID string, clientSecret string) string {
	credentials := clientID + ":" + clientSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}
