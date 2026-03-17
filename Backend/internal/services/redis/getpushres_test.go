package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	"github.com/go-redis/redismock/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func getMockParseResult() *models.ParseResult {
	return &models.ParseResult{
		UserID: 123,
		EntryData: []models.Entry{
			{
				UserID:          123,
				Title:           "Netflix",
				Price:           decimal.NewFromFloat(15.99),
				Currency:        "USD",
				Period:          models.Monthly,
				Category:        models.Entertainment,
				NextPaymentDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				IsActive:        true,
			},
		},
	}
}

func TestPushParseResult(t *testing.T) {
	db, mock := redismock.NewClientMock()
	conn := NewRedisConn(db)
	ctx := context.Background()

	result := getMockParseResult()
	payload, _ := json.Marshal(result)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectLPush(models.ResultQueueName, payload).SetVal(1)

		err := conn.PushParseResult(ctx, result)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Redis Error", func(t *testing.T) {
		mock.ExpectLPush(models.ResultQueueName, payload).SetErr(errors.New("redis connection lost"))

		err := conn.PushParseResult(ctx, result)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetParseResult(t *testing.T) {
	db, mock := redismock.NewClientMock()
	conn := NewRedisConn(db)
	ctx := context.Background()

	result := getMockParseResult()
	payload, _ := json.Marshal(result)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBRPop(0, models.ResultQueueName).SetVal([]string{models.ResultQueueName, string(payload)})

		retrievedResult, err := conn.GetParseResult(ctx)

		assert.NoError(t, err)
		assert.Equal(t, result.UserID, retrievedResult.UserID)
		assert.Equal(t, result.EntryData[0].Title, retrievedResult.EntryData[0].Title)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Redis Error", func(t *testing.T) {
		mock.ExpectBRPop(0, models.ResultQueueName).SetErr(errors.New("timeout or disconnect"))

		retrievedResult, err := conn.GetParseResult(ctx)

		assert.Error(t, err)
		assert.Nil(t, retrievedResult)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
