package entities

type Chat struct {
	SenderId   uint   `gorm:"sender_id"`
	ReceiverId uint   `gorm:"receiver_id"`
	Message    string `gorm:"message"`
}
