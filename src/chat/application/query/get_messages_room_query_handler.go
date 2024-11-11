package application

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type GetMessagesRoomQueryHandler struct {
	RedisClient *redis.Client
}

func (h *GetMessagesRoomQueryHandler) Handle(query *GetMessagesRoomQuery) (interface{}, error) {
	log.Printf("[GetMessagesRoomQuery] ->> Get message %s from room", query.Room)

	messages, err := h.RedisClient.XRead(context.Background(), &redis.XReadArgs{
		Streams: []string{query.Room, "0"},
		Block:   0,
	}).Result()

	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, msg := range messages[0].Messages {
		result = append(result, msg.Values)
	}

	return result, nil
}
