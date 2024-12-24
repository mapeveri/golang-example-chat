package application_test

import (
	application "distributed-chat/src/chat/application/query"
	"distributed-chat/src/chat/domain"
	"distributed-chat/src/chat/infrastructure/doubles"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessagesRoomQueryHandler_Handle(t *testing.T) {
	t.Run("should return an error when repository fails", func(t *testing.T) {
		mockRepo := new(doubles.MockMessageRepository)
		handler := application.GetMessagesRoomQueryHandler{
			MessageRepository: mockRepo,
		}
		query := &application.GetMessagesRoomQuery{Room: "test-room"}
		expectedError := errors.New("repository error")
		mockRepo.On("ByRoom", "test-room").Return(nil, expectedError)

		result, err := handler.Handle(query)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "ByRoom", "test-room")
	})

	t.Run("should return messages when repository succeeds", func(t *testing.T) {
		mockRepo := new(doubles.MockMessageRepository)
		handler := application.GetMessagesRoomQueryHandler{
			MessageRepository: mockRepo,
		}
		query := &application.GetMessagesRoomQuery{Room: "test-room"}
		expectedMessages := []domain.Message{
			{User: "test", Room: "test-room", Message: "Hello"},
			{User: "test2", Room: "test-room", Message: "World"},
		}
		mockRepo.On("ByRoom", "test-room").Return(expectedMessages, nil)

		result, err := handler.Handle(query)

		assert.NoError(t, err)
		assert.Equal(t, expectedMessages, result)
		mockRepo.AssertCalled(t, "ByRoom", "test-room")
	})
}
