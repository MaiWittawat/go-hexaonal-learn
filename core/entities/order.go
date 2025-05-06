package entities

type Order struct {
	UserId    uint `json:"user_id"`
	ProductId uint `json:"product_id"`
}
