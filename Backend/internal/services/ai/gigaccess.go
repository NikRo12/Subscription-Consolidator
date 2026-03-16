package ai

import (
	"errors"
	"os"
)

var ErrNoAuthKey = errors.New("There's no GigaChat auth-key")

/*
This function load variables from .env file to env-variables of the current process
It tries to get the gigachat API auth-key
If there're no AUTH_KEY variable, ir returns NoAuthKey error
*/
func GetGigaChatAuthKey() (string, error) {
	addr := os.Getenv("AUTH_KEY")
	if addr == "" {
		return "", ErrNoAuthKey
	}
	return addr, nil
}
