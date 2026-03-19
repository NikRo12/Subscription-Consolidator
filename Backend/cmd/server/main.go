package main

import (
	"flag"
	"log"

	"github.com/NikRo12/Subscription-Consolidator/Backend/configs"
	httpserver "github.com/NikRo12/Subscription-Consolidator/Backend/internal/transport/httpserver"
)

func main() {
	flag.Parse()

	if err := httpserver.Start(
		configs.GetDBURL(),
		configs.GetLogLevel(),
		configs.GetServerHost(),
		configs.GetRedisAddr(),
		configs.GetGoogleClientID(),
		configs.GetGoogleClientSecret(),
	); err != nil {
		log.Fatal(err)
	}

}
