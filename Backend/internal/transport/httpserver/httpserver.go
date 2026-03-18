package httpserver

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/email"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store/sqlstore"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/transport/queueconsumer"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func Start(databaseURL, logLevel, bindAddr, redisAddr, clientID, clientSecret string) error {
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

	authSerive := email.NewAuthService(clientID, clientSecret)

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer redisClient.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	consumer := queueconsumer.NewQueueConsumer(redisClient, store, logger)
	consumer.StartListeninig(ctx)

	s := newServer(store, logger, redisClient, authSerive)
	httpServer := &http.Server{
		Addr:    bindAddr,
		Handler: s,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	logger.Infof("Server started on %s", bindAddr)

	<-quit
	logger.Info("Shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	return httpServer.Shutdown(shutdownCtx)
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
