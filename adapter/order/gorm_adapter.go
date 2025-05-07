package adapter

import (
	"github.com/wittawat/go-hex/core/entities"
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) DeleteOne(id int) error {
	result := r.db.Delete(&entities.Order{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormOrderRepository) FindByUserId(userId int) ([]entities.Product, error) {
	var products []entities.Product
	result := r.db.Table("orders").
		Select("products.title, products.price, products.detail").
		Joins("JOIN products ON orders.product_id = products.id").
		Where("orders.user_id = ?", userId).
		Scan(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (r *GormOrderRepository) UpdateOne(order *entities.Order, id int) error {
	var o entities.Order
	result := r.db.First(&o, id)
	if result.Error != nil {
		return result.Error
	}
	o.UserId = order.UserId
	o.ProductId = order.ProductId

	result = r.db.Save(&o)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormOrderRepository) Save(order *entities.Order) error {
	result := r.db.Save(&order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
