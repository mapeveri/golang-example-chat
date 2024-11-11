package application

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type SendMessageCommandHandler struct {
	RedisClient *redis.Client
}

func (h *SendMessageCommandHandler) Handle(cmd *SendMessageCommand) error {
	log.Printf("[SendMessageCommandHandler] ->> Send message %s to room %s from %s", cmd.Message, cmd.Room, cmd.User)

	_, err := h.RedisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: cmd.Room,
		Values: map[string]interface{}{
			"user":    cmd.User,
			"message": cmd.Message,
		},
	}).Result()

	return err
}
