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
