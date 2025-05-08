package service

import (
	"github.com/wittawat/go-hex/core/entities"
	productPort "github.com/wittawat/go-hex/core/port/product"
)

type ProductService struct {
	repo productPort.ProductRepository
}

func NewProductService(repo productPort.ProductRepository) productPort.ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Save(product *entities.Product) error {
	if err := s.repo.Save(product); err != nil {
		return err
	}
	return nil
}

func (s *ProductService) Find() ([]entities.Product, error) {
	products, err := s.repo.Find()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) FindById(id int) (*entities.Product, error) {
	product, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateOne(product *entities.Product, id int) error {
	if err := s.repo.UpdateOne(product, id); err != nil {
		return err
	}
	return nil
}

func (s *ProductService) DeleteOne(id int) error {
	if err := s.repo.DeleteOne(id); err != nil {
		return err
	}
	return nil
}
