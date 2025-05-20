package validator

import (
	"time"

	"github.com/wittawat/go-hex/core/entities"
)

func updateOrdedr(oldOrder *entities.Order, newOrder *entities.Order) *entities.Order {
	if newOrder.ProductID == "" {
		newOrder.ProductID = oldOrder.ProductID
	}
	newOrder.UserID = oldOrder.UserID
	newOrder.CreatedAt = oldOrder.CreatedAt
	newOrder.UpdatedAt = time.Now()
	return newOrder
}

func EnsureUpdateOrder(oldOrder *entities.Order, newOrder *entities.Order) *entities.Order {
	order := updateOrdedr(oldOrder, newOrder)
	return order
}
