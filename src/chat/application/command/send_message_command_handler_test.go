package application_test

import (
	application "distributed-chat/src/chat/application/command"
	"distributed-chat/src/chat/domain"
	"distributed-chat/src/chat/infrastructure/doubles"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessageCommandHandler_Handle(t *testing.T) {
	t.Run("should send a message to a room", func(t *testing.T) {
		mockRepo := new(doubles.MockMessageRepository)
		handler := application.SendMessageCommandHandler{MessageRepository: mockRepo}
		cmd := &application.SendMessageCommand{
			User:    "user1",
			Message: "Hello, World!",
			Room:    "room1",
		}
		expectedMessage := domain.Message{
			User:    "user1",
			Message: "Hello, World!",
			Room:    "room1",
		}

		mockRepo.On("Save", expectedMessage).Return(nil)

		err := handler.Handle(cmd)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Save", expectedMessage)

	})
}
