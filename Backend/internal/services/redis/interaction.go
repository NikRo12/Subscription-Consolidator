package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	"github.com/redis/go-redis/v9"
)

/*
Wrapping aroung redis-client
*/
type RedisConn struct {
	client *redis.Client
}

func NewRedisConn(client *redis.Client) *RedisConn {
	return &RedisConn{client: client}
}

/*
Closing connection
*/
func (conn *RedisConn) terminate() error {
	if err := conn.client.Close(); err != nil {
		return fmt.Errorf("[ERROR]! Redis-client: %v", err)
	}
	return nil
}

/*
Push the task to a redis list
*/
func (conn *RedisConn) PushTask(ctx context.Context, task *models.Task) error {
	byteData, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("[ERROR]! Json-marhslling: %v", err)
	}
	err = conn.client.LPush(ctx, models.TaskQueueName, byteData).Err()
	if err != nil {
		return fmt.Errorf("[ERROR]! Redis-lpush: %v", err)
	}
	return nil
}

/*
Blocking operation!
It waits for any tasks in the queue and takes it
*/
func (conn *RedisConn) GetTask(ctx context.Context) (*models.Task, error) {
	task := &models.Task{}
	res, err := conn.client.BRPop(ctx, 0, models.TaskQueueName).Result()
	if err != nil {
		return nil, fmt.Errorf("[ERROR]! Redis-brpop: %v", err)
	}
	byteData := []byte(res[1])
	err = json.Unmarshal(byteData, task)
	if err != nil {
		return nil, fmt.Errorf("[ERROR]! Json-marhslling: %v", err)
	}
	return task, nil
}
