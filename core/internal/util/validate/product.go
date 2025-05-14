package validator

import (
	"errors"

	"github.com/wittawat/go-hex/core/entities"
)

func updateProduct(oldProduct *entities.Product, newProduct *entities.Product) *entities.Product {
	if newProduct.Title == "" {
		newProduct.Title = oldProduct.Title
	}
	if newProduct.Price == 0 {
		newProduct.Price = oldProduct.Price
	}
	if newProduct.Detail == "" {
		newProduct.Detail = oldProduct.Detail
	}
	newProduct.CreatedBy = oldProduct.CreatedBy
	return newProduct
}

func IsValidProduct(product *entities.Product) error {
	if len(product.Title) < 4 {
		return errors.New("title must more than 4 character")
	}
	if len(product.Detail) < 12 {
		return errors.New("detail must more than 12 character")
	}
	return nil
}

func EnsureUpdateProduct(oldProduct *entities.Product, newProduct *entities.Product) (*entities.Product, error) {
	product := updateProduct(oldProduct, newProduct)
	if err := IsValidProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}
