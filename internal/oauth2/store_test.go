package oauth2

import (
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
	assert.Nil(t, err)
	defer func() {
		err := os.RemoveAll(notionDir)
		assert.Nil(t, err)
	}()

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
