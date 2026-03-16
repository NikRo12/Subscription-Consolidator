package redis

import (
	"os"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/commonerrors"
	"github.com/joho/godotenv"
)

/*
This function load variables from .env file to env-variables of the current process
It tries to get the redis' addr
If there're no REDIS_ADDR variable, it returns `localhost:6379`
*/
func GetRedisAddr() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", commonerrors.NoENVFile
	}
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return addr, nil
}
