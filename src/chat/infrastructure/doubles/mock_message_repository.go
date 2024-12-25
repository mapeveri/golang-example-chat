package doubles

import (
	"distributed-chat/src/chat/domain"

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
