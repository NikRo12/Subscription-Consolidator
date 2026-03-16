package main

import (
	"flag"
	"log"

	"github.com/NikRo12/Subscription-Consolidator/Backend/configs"
	httpserver "github.com/NikRo12/Subscription-Consolidator/Backend/internal/transport/httpserver"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("/home/nikita081105/develope/Subscription-Consolidator/Backend/.env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if err := httpserver.Start(configs.GetDBDriver(), configs.GetDBURL(), configs.GetLogLevel(), configs.GetServerHost()); err != nil {
		log.Fatal(err)
	}

}
