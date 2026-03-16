package configs

import "os"

func GetServerHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost:8080"
	}
	return host
}

func GetDBDriver() string {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "postgres"
	}
	return driver
}

func GetDBURL() string {
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
