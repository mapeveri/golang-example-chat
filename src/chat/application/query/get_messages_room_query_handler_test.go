package application_test

import (
	application "distributed-chat/src/chat/application/query"
	"distributed-chat/src/chat/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) ByRoom(room string) (interface{}, error) {
	args := m.Called(room)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) Save(message domain.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestGetMessagesRoomQueryHandler_Handle(t *testing.T) {
	t.Run("should return an error when repository fails", func(t *testing.T) {
		mockRepo := new(MockMessageRepository)
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
		mockRepo := new(MockMessageRepository)
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
