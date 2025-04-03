package oauth2

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldAuth(t *testing.T) {
	originalDir := notionDir
	defer func() { notionDir = originalDir }()

	var err error
	notionDir, err = os.MkdirTemp("", ".notion-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(notionDir)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Test without token
	auth := shouldAuth()
	assert.True(t, auth)

	// Test with token
	err = storeToken("{\"access_token\": \"test-token\"}")
	assert.Nil(t, err)

	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"y", true},
		{"Y", true},
		{"n", false},
		{"N", false},
		{"test", false},
	}

	for _, test := range tests {
		scanln = func(a ...any) (n int, err error) {
			assert.Len(t, a, 1)

			arg := a[0]
			*arg.(*string) = test.input
			return 0, nil
		}
		auth = shouldAuth()
		assert.Equal(t, test.want, auth)
	}

	scanln = func(a ...any) (n int, err error) {
		return 0, errors.New("err")
	}
	auth = shouldAuth()
	assert.False(t, auth)

	defer func() { scanln = fmt.Scanln }()
}

func TestValidateCodeResponse(t *testing.T) {
	tests := []struct {
		res   codeResponse
		state string
		want  bool
	}{
		{codeResponse{}, "", false},
		{codeResponse{}, "state", false},
		{codeResponse{state: "state"}, "", false},
		{codeResponse{state: "state"}, "state", true},
		{codeResponse{err: "err"}, "", false},
	}

	for _, test := range tests {
		valid := validateCodeResponse(test.res, test.state)
		assert.Equal(t, test.want, valid)
	}
}
