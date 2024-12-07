package repositories

import (
	"context"
	"distributed-chat/src/chat/domain"

	"github.com/go-redis/redis/v8"
)

type RedisMessageRepository struct {
	RedisClient *redis.Client
}

func NewRedisMessageRepository(redisClient *redis.Client) *RedisMessageRepository {
	return &RedisMessageRepository{
		RedisClient: redisClient,
	}
}

func (r *RedisMessageRepository) Save(message domain.Message) error {
	_, err := r.RedisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: message.Room,
		Values: map[string]interface{}{
			"user":    message.User,
			"message": message.Message,
		},
	}).Result()

	return err
}

func (r *RedisMessageRepository) ByRoom(room string) (interface{}, error) {
	messages, err := r.RedisClient.XRead(context.Background(), &redis.XReadArgs{
		Streams: []string{room, "0"},
		Block:   0,
	}).Result()

	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, msg := range messages[0].Messages {
		result = append(result, msg.Values)
	}

	return result, err
}
