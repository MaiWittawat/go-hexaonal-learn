package request

type OrderRequest struct {
	UserId    uint `json:"user_id"`
	ProductId uint `json:"product_id"`
}
