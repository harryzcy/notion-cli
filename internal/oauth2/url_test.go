package oauth2

import (
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildAuthURL(t *testing.T) {
	tests := []struct {
		clientID string
		state    string
	}{
		{
			clientID: "12345",
			state:    "abcde",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := buildAuthURL(tt.clientID, tt.state)

			urlObject, err := url.Parse(got)
			assert.Nil(t, err)
			assert.Equal(t, "https", urlObject.Scheme)
			assert.Equal(t, "api.notion.com", urlObject.Host)
			assert.Equal(t, "/v1/oauth/authorize", urlObject.Path)
			assert.Len(t, urlObject.Query(), 5)
			assert.Equal(t, tt.clientID, urlObject.Query().Get("client_id"))
			assert.Equal(t, "http://localhost:"+port+"/callback", urlObject.Query().Get("redirect_uri"))
			assert.Equal(t, "code", urlObject.Query().Get("response_type"))
			assert.Equal(t, "user", urlObject.Query().Get("owner"))
			assert.Equal(t, tt.state, urlObject.Query().Get("state"))
		})
	}
}

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		params map[string]string
		want   string
	}{
		{
			params: map[string]string{
				"foo": "bar",
			},
		},
		{
			params: map[string]string{
				"foo":  "bar",
				"foo2": "bar2",
			},
		},
	}

	for _, tt := range tests {
		got := buildQuery(tt.params)

		assert.Equal(t, len(tt.params)-1, strings.Count(got, "&"))
		assert.False(t, strings.HasPrefix(got, "&"))
		assert.False(t, strings.HasSuffix(got, "&"))
		for key, value := range tt.params {
			assert.Contains(t, got, key+"="+value)
		}
	}
}
