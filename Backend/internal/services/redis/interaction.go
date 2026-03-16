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
func (conn *RedisConn) Terminate() error {

	if err := conn.client.Close(); err != nil {
		return fmt.Errorf("[ERROR]! redis-client: %w", err)
	}

	return nil
}

/*
Push the task to a redis list
*/
func (conn *RedisConn) PushTask(
	ctx context.Context,
	task *models.Task) error {

	byteData, err := json.Marshal(task)

	if err != nil {
		return fmt.Errorf("[ERROR]! json-marhslling: %w", err)
	}

	err = conn.client.LPush(ctx,
		models.TaskQueueName,
		byteData).Err()

	if err != nil {
		return fmt.Errorf("[ERROR]! redis-lpush: %w", err)
	}

	return nil
}

/*
Blocking operation!
It waits for any tasks in the queue and takes it
*/
func (conn *RedisConn) GetTask(ctx context.Context) (
	*models.Task,
	error) {
	task := &models.Task{}

	res, err := conn.client.BRPop(ctx, 0, models.TaskQueueName).Result()

	if err != nil {
		return nil, fmt.Errorf("[ERROR]! redis-brpop: %w", err)
	}

	if len(res) != 2 {
		return nil, fmt.Errorf("unexpected redis-format")
	}

	byteData := []byte(res[1])

	err = json.Unmarshal(byteData, task)

	if err != nil {
		return nil, fmt.Errorf("[ERROR]! json-marhslling: %w", err)
	}

	return task, nil
}

func (conn *RedisConn) PushParseResult(
	ctx context.Context,
	parseRes *models.ParseResult) error {

	serializedData, err := json.Marshal(parseRes)
	if err != nil {
		return fmt.Errorf("[ERROR]! json marshalling: %w", err)
	}

	err = conn.client.LPush(ctx,
		models.ResultQueueName,
		serializedData).Err()

	if err != nil {
		return fmt.Errorf("[ERROR]! сannot push the result: %w", err)
	}

	return nil
}

func (conn *RedisConn) GetParseResult(ctx context.Context) (
	*models.ParseResult,
	error) {
	parseRes := &models.ParseResult{}

	redisData, err := conn.client.BRPop(ctx, 0, models.ResultQueueName).Result()

	if err != nil {
		return nil, fmt.Errorf("[ERROR]! сannot read result: %w", err)
	}

	if len(redisData) != 2 {
		return nil, fmt.Errorf("unexpected redis-format")
	}

	serializedRes := []byte(redisData[1])

	err = json.Unmarshal(serializedRes, parseRes)

	if err != nil {
		return nil, fmt.Errorf("[ERROR]! сannot unmarshal result: %w", err)
	}

	return parseRes, nil
}
