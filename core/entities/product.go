package entities

type Product struct {
	Title  string `json:"title"`
	Price  uint   `json:"price"`
	Detail string `json:"detail"`
}
