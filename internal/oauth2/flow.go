package oauth2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/browser"
)

var (
	baseURL      = "https://api.notion.com/v1/oauth/authorize"
	clientID     = os.Getenv("NOTION_CLIENT_ID")
	clientSecret = os.Getenv("NOTION_CLIENT_SECRET")
	port         = "30001"
)

func Flow() {
	toAuth := shouldAuth()
	if !toAuth {
		return
	}

	c := make(chan codeResponse)
	go startServer(c)

	state := generateRandomState()

	url := buildAuthURL(state)
	_ = browser.OpenURL(url)

	res := <-c
	valid := validateCodeResponse(res, state)
	if !valid {
		return
	}

	err := exchangeGrant(res.code, state)
	if err != nil {
		fmt.Println("Error exchanging grant:", err)
		return
	}

}

func validateCodeResponse(res codeResponse, state string) bool {
	if res.err != "" {
		fmt.Println("Access is not granted")
		return false
	}

	if res.state != state {
		fmt.Println("State is not matched")
		return false
	}
	return true
}

func exchangeGrant(code string, state string) error {
	url := "https://api.notion.com/v1/oauth/token"

	data := map[string]string{
		"grant_type":   "authorization_code",
		"code":         code,
		"redirect_uri": "http://localhost:" + port + "/callback",
	}
	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", basicAuth(clientID, clientSecret))

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	var token tokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		return err
	}

	err = storeToken(string(body))

	return err
}

func basicAuth(clientID string, clientSecret string) string {
	credentials := clientID + ":" + clientSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}

func shouldAuth() bool {
	token, err := getToken()
	if err == nil && token.AccessToken != "" {
		// code is stored, ask if user wants to re-auth
		fmt.Print("Do you want to re-auth? (y/N): ")
		var input string
		fmt.Scanln(&input)
		return strings.ToLower(input) == "y"
	}
	return true
}
