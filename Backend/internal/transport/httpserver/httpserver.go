package httpserver

import (
	"database/sql"
	"net/http"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store/sqlstore"
	"github.com/sirupsen/logrus"
)

func Start(databaseURL, logLevel, bindAddr string) error {
	db, err := newDB(databaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.NewSqlStore(db)

	logger, err := configureLogger(logLevel)
	if err != nil {
		return err
	}
	s := newServer(store, logger)

	return http.ListenAndServe(bindAddr, s)
}

func newDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func configureLogger(logLevel string) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(level)

	return logger, nil
}
