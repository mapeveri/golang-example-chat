package domain

type MessageRepository interface {
	Save(message Message) error

	ByRoom(room string) (interface{}, error)
}
