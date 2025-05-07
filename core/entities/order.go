package entities

type Order struct {
	ID        uint
	UserId    uint `json:"user_id"`
	ProductId uint `json:"product_id"`
}
