package main

import (
	"context"
	commands "distributed-chat/src/chat/application/command"
	queries "distributed-chat/src/chat/application/query"
	"distributed-chat/src/chat/domain"
	"distributed-chat/src/shared/infrastructure/bus"
	"distributed-chat/src/shared/infrastructure/persistence/redis/repositories"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/dig"
)

func buildContainer() *dig.Container {
	container := dig.New()

	rdb := redisConnection()
	container.Provide(func() *redis.Client {
		return rdb
	})

	container.Provide(func(rdb *redis.Client) domain.MessageRepository {
		return repositories.NewRedisMessageRepository(rdb)
	})

	container.Provide(func(messageRepository domain.MessageRepository) *commands.SendMessageCommandHandler {
		return &commands.SendMessageCommandHandler{
			MessageRepository: messageRepository,
		}
	})

	container.Provide(func(rdb *redis.Client) *queries.GetMessagesRoomQueryHandler {
		return &queries.GetMessagesRoomQueryHandler{
			RedisClient: rdb,
		}
	})

	return container
}

func main() {
	container := buildContainer()

	commandBus := bus.NewMemoryCommandBus()

	commandBus.Register(&commands.SendMessageCommand{}, func(cmd interface{}) error {
		var handler *commands.SendMessageCommandHandler
		if err := container.Invoke(func(h *commands.SendMessageCommandHandler) {
			handler = h
		}); err != nil {
			log.Fatalf("Error while fetching handler: %v", err)
		}

		return handler.Handle(cmd.(*commands.SendMessageCommand))
	})

	queryBus := bus.NewMemoryQueryBus()

	queryBus.Register(&queries.GetMessagesRoomQuery{}, func(cmd interface{}) (interface{}, error) {
		var handler *queries.GetMessagesRoomQueryHandler
		if err := container.Invoke(func(h *queries.GetMessagesRoomQueryHandler) {
			handler = h
		}); err != nil {
			log.Fatalf("Error while fetching handler: %v", err)
		}

		return handler.Handle(cmd.(*queries.GetMessagesRoomQuery))
	})

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

		sendMessageCommand := &commands.SendMessageCommand{User: request.User, Message: request.Message, Room: request.Room}
		err := commandBus.Execute(sendMessageCommand)
		if err != nil {
			log.Printf("Failed to send message: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
	})

	r.GET("/api/receive-messages", func(c *gin.Context) {
		room := c.DefaultQuery("room", "default-room")

		getMessagesRoomQuery := &queries.GetMessagesRoomQuery{Room: room}
		messages, err := queryBus.Execute(getMessagesRoomQuery)
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

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error trying to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")

	return rdb
}
