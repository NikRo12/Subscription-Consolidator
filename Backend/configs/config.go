package configs

import "os"

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
