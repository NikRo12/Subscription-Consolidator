package ai

import (
	"errors"
	"os"
)

var ErrNoAuthKey = errors.New("There's no GigaChat auth-key")

func getGigaChatAuthKey() (string, error) {
	addr := os.Getenv("AUTH_KEY")
	if addr == "" {
		return "", ErrNoAuthKey
	}
	return addr, nil
}
