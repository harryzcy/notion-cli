package oauth2

import (
	"math/rand"
	"time"
)

const stateLength = 5

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomState() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, stateLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
