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

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

var rdb = redisConnection()

func main() {
	r := gin.Default()

	r.POST("/send-message", sendMessageHandler)
	r.GET("/receive-messages", receiveMessagesHandler)

	log.Println("Server started on :8080")
	r.Run(":8080")
}

func redisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")

	return rdb
}

func sendMessageHandler(c *gin.Context) {
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	room := c.DefaultQuery("room", "")
	if room == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room parameter is required"})
		return
	}

	err := sendMessage(rdb, room, msg.User, msg.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
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

func receiveMessagesHandler(c *gin.Context) {
	room := c.DefaultQuery("room", "")
	if room == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room parameter is required"})
		return
	}

	lastID := c.DefaultQuery("last_id", "0")

	messages, err := receiveMessages(rdb, room, lastID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error receiving messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func receiveMessages(rdb *redis.Client, room string, lastID string) ([]Message, error) {
	var messages []Message

	messagesRedis, err := rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{room, lastID},
		Block:   0,
	}).Result()

	if err != nil {
		return nil, err
	}

	for _, msg := range messagesRedis[0].Messages {
		user := msg.Values["user"].(string)
		message := msg.Values["message"].(string)
		messages = append(messages, Message{
			User:    user,
			Message: message,
		})
		lastID = msg.ID
	}

	return messages, nil
}
