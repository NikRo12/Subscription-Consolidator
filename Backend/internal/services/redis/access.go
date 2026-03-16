package redis

import (
	"os"
)

/*
This function load variables from .env file to env-variables of the current process
It tries to get the redis' addr
If there're no REDIS_ADDR variable, it returns `localhost:6379`
*/
func GetRedisAddr() (string, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return addr, nil
}
