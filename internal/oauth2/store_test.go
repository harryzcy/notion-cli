package oauth2

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotionDir(t *testing.T) {
	assert.True(t, strings.HasSuffix(notionDir, "/.notion"))
}

func TestToken(t *testing.T) {
	originalDir := notionDir
	defer func() { notionDir = originalDir }()

	var err error
	notionDir, err = os.MkdirTemp("", ".notion-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(notionDir)

	// Test that the token is not found
	token, err := GetToken()
	assert.NotNil(t, err)
	assert.Empty(t, token)

	// Test that the token is stored
	err = storeToken("{\"access_token\": \"test-token\"}")
	assert.Nil(t, err)

	// Test that the token is retrieved
	token, err = GetToken()
	assert.Nil(t, err)
	assert.Equal(t, "test-token", token.AccessToken)
}

func TestStoreToken_PermissionDenied(t *testing.T) {
	originalDir := notionDir
	defer func() { notionDir = originalDir }()

	var err error
	notionDir, err = os.MkdirTemp("", ".notion-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(notionDir)

	err = os.Chmod(notionDir, 0444) // r--r--r--
	assert.Nil(t, err)
	notionDir += "/.notion" // a directory inside of a read-only directory

	err = storeToken("{\"access_token\": \"test-token\"}")
	assert.NotNil(t, err)
}
