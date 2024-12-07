package application

import (
	"distributed-chat/src/chat/domain"
	"log"
)

type SendMessageCommandHandler struct {
	MessageRepository domain.MessageRepository
}

func (h *SendMessageCommandHandler) Handle(cmd *SendMessageCommand) error {
	log.Printf("[SendMessageCommandHandler] ->> Send message %s to room %s from %s", cmd.Message, cmd.Room, cmd.User)

	message := domain.Message{
		User:    cmd.User,
		Message: cmd.Message,
		Room:    cmd.Room,
	}

	err := h.MessageRepository.Save(message)

	return err
}
