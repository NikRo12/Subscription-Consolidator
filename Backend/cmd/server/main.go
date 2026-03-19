package main

import (
	"flag"
	"log"

	"github.com/NikRo12/Subscription-Consolidator/Backend/configs"
	httpserver "github.com/NikRo12/Subscription-Consolidator/Backend/internal/transport/httpserver"
)

func main() {
	flag.Parse()

	clientID, err := configs.GetGoogleClientID()
	if err != nil {
		log.Fatal(err)
		return
	}

	clientSecret, err := configs.GetGoogleClientSecret()
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := httpserver.Start(
		configs.GetDBURL(),
		configs.GetLogLevel(),
		configs.GetServerHost(),
		configs.GetRedisAddr(),
		clientID,
		clientSecret,
	); err != nil {
		log.Fatal(err)
	}

}
