package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomState(t *testing.T) {
	state := generateRandomState()
	assert.NotEmpty(t, state)
	assert.Len(t, state, stateLength)
}
