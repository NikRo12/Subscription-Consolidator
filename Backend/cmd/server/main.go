package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	httpserver "github.com/NikRo12/Subscription-Consolidator/Backend/internal/transport/httpserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "../../configs/httpserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := httpserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	server := httpserver.New(config)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
