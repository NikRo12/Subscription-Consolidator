package store_test

import (
	"os"
	"testing"
)

var (
	driver      string
	databaseURL string
)

func TestMain(m *testing.M) {
	driver = os.Getenv("DATABASE_DRIVER")
	databaseURL = os.Getenv("DATABASE_URL")

	if driver == "" {
		driver = "postgres"
	}

	if databaseURL == "" {
		databaseURL = "host=localhost port=5432 user=nikita081105 password=Elfhybr081105 dbname=subconapp_test sslmode=disable"
	}

	os.Exit(m.Run())
}
