package adapter

import (
	"github.com/wittawat/go-hex/core/entities"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

func (r *GormProductRepository) Save(product *entities.Product) error {
	result := r.db.Create(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormProductRepository) Find() ([]entities.Product, error) {
	var products []entities.Product
	result := r.db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (r *GormProductRepository) FindById(id int) (*entities.Product, error) {
	var product entities.Product
	result := r.db.First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *GormProductRepository) UpdateOne(product *entities.Product, id int) error {
	result := r.db.Model(&entities.Product{}).Where("id = ?", id).
		Select("title", "price", "description").
		Updates(product)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormProductRepository) DeleteOne(id int) error {
	result := r.db.Delete(id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
