package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestPushTask(t *testing.T) {
	db, mock := redismock.NewClientMock()
	conn := NewRedisConn(db)
	ctx := context.Background()

	task := &models.Task{TaskID: "1", UserID: "user_123"}
	payload, _ := json.Marshal(task)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectLPush(models.TaskQueueName, payload).SetVal(1)

		err := conn.PushTask(ctx, task)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Redis Error", func(t *testing.T) {
		mock.ExpectLPush(models.TaskQueueName, payload).SetErr(errors.New("redis connection lost"))

		err := conn.PushTask(ctx, task)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Redis-lpush")
	})
}

func TestGetTask(t *testing.T) {
	db, mock := redismock.NewClientMock()
	conn := NewRedisConn(db)
	ctx := context.Background()

	task := &models.Task{TaskID: "1", UserID: "user_123"}
	payload, _ := json.Marshal(task)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBRPop(0, models.TaskQueueName).SetVal([]string{models.TaskQueueName, string(payload)})

		result, err := conn.GetTask(ctx)

		assert.NoError(t, err)
		assert.Equal(t, task.TaskID, result.TaskID)
		assert.Equal(t, task.UserID, result.UserID)
	})

	t.Run("Unmarshal Error", func(t *testing.T) {
		mock.ExpectBRPop(0, models.TaskQueueName).SetVal([]string{models.TaskQueueName, "invalid-json"})

		result, err := conn.GetTask(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Json-marhslling")
	})
}
