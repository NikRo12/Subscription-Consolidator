package main

import (
	"flag"
	"log"

	"github.com/NikRo12/Subscription-Consolidator/Backend/configs"
	httpserver "github.com/NikRo12/Subscription-Consolidator/Backend/internal/transport/httpserver"
)

var (
	serverConfigPath  string
	storageConfigPath string
)

func init() {
	flag.StringVar(&serverConfigPath, "server_config-path", "../configs/server.toml", "path to server config file")
	flag.StringVar(&storageConfigPath, "storage_config-path", "../configs/storage.toml", "path to storage config file")
}

func main() {
	flag.Parse()

	srvConfig := configs.LoadServerConfig(serverConfigPath)
	strConfig := configs.LoadStorageConfig(storageConfigPath)

	if err := httpserver.Start(srvConfig, strConfig); err != nil {
		log.Fatal(err)
	}

}
