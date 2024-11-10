package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redisConnection()

	r := gin.Default()

	r.StaticFile("/", "./index.html")

	r.POST("/api/send-message", func(c *gin.Context) {
		var request struct {
			User    string `json:"user"`
			Message string `json:"message"`
			Room    string `json:"room"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err := sendMessage(rdb, request.Room, request.User, request.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
	})

	r.GET("/api/receive-messages", func(c *gin.Context) {
		room := c.DefaultQuery("room", "default-room")

		messages, err := receiveMessages(rdb, room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to receive messages"})
			return
		}

		c.JSON(http.StatusOK, messages)
	})

	r.Run(":8080")
}

func redisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error trying to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")

	return rdb
}

func sendMessage(rdb *redis.Client, room string, user string, message string) error {
	_, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: room,
		Values: map[string]interface{}{
			"user":    user,
			"message": message,
		},
	}).Result()

	return err
}

func receiveMessages(rdb *redis.Client, room string) ([]map[string]interface{}, error) {
	messages, err := rdb.XRead(ctx, &redis.XReadArgs{
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

	return result, nil
}
