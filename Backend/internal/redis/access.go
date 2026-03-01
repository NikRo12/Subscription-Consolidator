package redis

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

/*
This error says that the root directory doesn't have .env-file
*/
var NoENVFile = errors.New("Cannot find .env file in project's root directory")

/*
This function load variables from .env file to env-variables of the current process
It tries to get the redis' addr
If there're no REDIS_ADDR variable, it returns `localhost:6379`
*/
func GetRedisAddr() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", NoENVFile
	}
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return addr, nil
}
