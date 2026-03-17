package main

import (
	"flag"
	"log"

	"github.com/NikRo12/Subscription-Consolidator/Backend/configs"
	httpserver "github.com/NikRo12/Subscription-Consolidator/Backend/internal/transport/httpserver"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if err := httpserver.Start(configs.GetDBURL(), configs.GetLogLevel(), configs.GetServerHost(), configs.GetRedisAddr()); err != nil {
		log.Fatal(err)
	}

}
