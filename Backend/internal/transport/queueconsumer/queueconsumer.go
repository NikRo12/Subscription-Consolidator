package queueconsumer

import (
	"context"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	intredis "github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/redis"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type QueueConsumer struct {
	redisConnection *intredis.RedisConn
	data            chan models.ParseResult
	store           store.Store
	logger          *logrus.Logger
}

func NewQueueConsumer(client *redis.Client, store store.Store, logger *logrus.Logger) *QueueConsumer {
	return &QueueConsumer{
		redisConnection: intredis.NewRedisConn(client),
		data:            make(chan models.ParseResult),
		store:           store,
		logger:          logger,
	}
}

func (qc *QueueConsumer) StartListeninig(ctx context.Context) {
	go qc.listen(ctx)
	go qc.consume(ctx)
}

func (qc *QueueConsumer) ParseResult(res models.ParseResult) {
	userID := res.UserID

	if err := qc.store.Sub().DeleteAllUserSubs(userID); err != nil {
		qc.logger.Errorf("failed to clear user_subs for user %d: %v", userID, err)
		return
	}

	for _, entry := range res.EntryData {
		sub := models.Subscription{
			Title:       entry.Title,
			Currency:    entry.Currency,
			Category:    entry.Category,
			IconURL:     entry.IconURL,
			BrandColor:  entry.BrandColor,
			Description: entry.Description,
		}

		if err := qc.store.Sub().CreateSub(&sub); err != nil {
			qc.logger.Error(err)
			continue
		}

		userSub := models.UserSubscription{
			UserID:          userID,
			SubID:           sub.ID,
			Price:           entry.Price,
			Period:          entry.Period,
			NextPaymentDate: entry.NextPaymentDate,
			IsActive:        entry.IsActive,
		}

		if err := qc.store.Sub().CreateUserSub(&userSub); err != nil {
			qc.logger.Error(err)
		}
	}
}

func (qc *QueueConsumer) listen(ctx context.Context) {
	for {
		res, err := qc.redisConnection.GetParseResult(ctx)
		if err != nil {
			qc.logger.Error(err)
			if ctx.Err() != nil {
				return
			}

			continue
		}

		select {
		case <-ctx.Done():
			return
		case qc.data <- *res:
		}
	}
}

func (qc *QueueConsumer) consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-qc.data:
			qc.ParseResult(res)
		}
	}
}
