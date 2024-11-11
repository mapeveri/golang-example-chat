package application

type SendMessageCommand struct {
	User    string
	Message string
	Room    string
}
