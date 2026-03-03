package httpserver

import (
	"database/sql"
	"net/http"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store/sqlstore"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseDriver, config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	s := newServer(store)

	return http.ListenAndServe(config.BindAddr, s)
}

func newDB(driver, url string) (*sql.DB, error) {
	db, err := sql.Open(driver, url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
