package task

import (
	"context"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	intredis "github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/redis"
	"github.com/redis/go-redis/v9"
)

type TaskService struct {
	redisConnection *intredis.RedisConn
}

func NewTaskService(client redis.Client) *TaskService {
	return &TaskService{
		redisConnection: intredis.NewRedisConn(&client),
	}
}

func (ts *TaskService) SendTask(ctx context.Context, task *models.Task) error {
	return ts.redisConnection.PushTask(ctx, task)
}
