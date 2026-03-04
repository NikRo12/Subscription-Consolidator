package configs

import (
	"log"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
}

type StorageConfig struct {
	DatabaseDriver string `toml:"database_driver"`
	DatabaseURL    string `toml:"database_url"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

func NewStorageConfig() *StorageConfig {
	return &StorageConfig{
		DatabaseDriver: "postgres",
	}
}

func LoadServerConfig(configPath string) *ServerConfig {
	var srvconfig ServerConfig

	_, err := toml.DecodeFile(configPath, &srvconfig)
	if err != nil {
		log.Fatal(err)
	}

	return &srvconfig
}

func LoadStorageConfig(configPath string) *StorageConfig {
	var strconfig StorageConfig

	_, err := toml.DecodeFile(configPath, &strconfig)
	if err != nil {
		log.Fatal(err)
	}

	return &strconfig
}
