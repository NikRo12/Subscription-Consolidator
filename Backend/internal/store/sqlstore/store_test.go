package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseDriver string
	databaseURL    string
)

func TestMain(m *testing.M) {
	databaseDriver = os.Getenv("DATABASE_DRIVER")
	if databaseDriver == "" {
		databaseDriver = "postgres"
	}

	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = `host=localhost port=5432 
			user=nikita081105 
			password=Elfhybr081105 
			dbname=subconapp_test 
			sslmode=disable`
	}

	os.Exit(m.Run())
}
