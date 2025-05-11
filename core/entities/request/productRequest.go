package request

type ProductRequest struct {
	Title  string `json:"title"`
	Price  uint   `json:"price"`
	Detail string `json:"detail"`
}
