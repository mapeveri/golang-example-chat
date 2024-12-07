package application

import (
	"distributed-chat/src/chat/domain"
	"log"
)

type GetMessagesRoomQueryHandler struct {
	MessageRepository domain.MessageRepository
}

func (h *GetMessagesRoomQueryHandler) Handle(query *GetMessagesRoomQuery) (interface{}, error) {
	log.Printf("[GetMessagesRoomQuery] ->> Get message %s from room", query.Room)

	messages, err := h.MessageRepository.ByRoom(query.Room)

	return messages, err
}
