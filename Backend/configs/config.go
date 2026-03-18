package configs

import (
	"errors"
	"os"
)

func GetServerHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost:8080"
	}
	return host
}

func GetDBURL() string {
	//TODO add empty value handler
	url := os.Getenv("DB_URL")
	return url
}

func GetTestDBURL() string {
	//TODO add empty value handler
	url := os.Getenv("TEST_DB_URL")
	return url
}

func GetLogLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "DEBUG"
	}
	return level
}

func GetJWTSecret() string {
	//TODO add empty value handler
	secret := os.Getenv("JWT_SECRET")
	return secret
}

func GetRedisAddr() string {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return addr
}

func GetGoogleClientID() (string, error) {
	cleintID := os.Getenv("CLIENT_ID")
	if cleintID == "" {
		return "", errors.New("cannot get client-id")
	}
	return cleintID, nil
}

func GetGoogleClientSecret() (string, error) {
	cleintSecret := os.Getenv("CLIENT_SECRET")
	if cleintSecret == "" {
		return "", errors.New("cannot get client-id")
	}
	return cleintSecret, nil
}
