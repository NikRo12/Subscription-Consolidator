package httpserver

import (
	"database/sql"
	"net/http"

	"github.com/NikRo12/Subscription-Consolidator/Backend/configs"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store/sqlstore"
	"github.com/sirupsen/logrus"
)

func Start(srvConfig *configs.ServerConfig, strConfig *configs.StorageConfig) error {
	db, err := newDB(strConfig.DatabaseDriver, strConfig.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)

	logger, err := configureLogger(srvConfig.LogLevel)
	if err != nil {
		return err
	}
	s := newServer(store, logger)

	return http.ListenAndServe(srvConfig.BindAddr, s)
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

func configureLogger(logLevel string) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(level)

	return logger, nil
}
