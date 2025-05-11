package chat

type ChatService interface {
	Create(message string) error
	Delete(id int) error
}
