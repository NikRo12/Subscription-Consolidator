package sqlstore_test

import (
	"log"
	"os"
	"testing"

	"github.com/NikRo12/Subscription-Consolidator/Backend/configs"
	"github.com/joho/godotenv"
)

var (
	databaseDriver string
	databaseURL    string
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("/home/nikita081105/develope/Subscription-Consolidator/Backend/.env"); err != nil {
		log.Fatal(err)
	}

	databaseDriver = "postgres"

	databaseURL = configs.GetTestDBURL()

	os.Exit(m.Run())
}
