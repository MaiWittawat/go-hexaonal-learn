package service

import (
	"github.com/wittawat/go-hex/core/entities"
	port "github.com/wittawat/go-hex/core/port/product"
)

type ProductService struct {
	ob port.ProductOutbound
}

func NewProductService(ob port.ProductOutbound) port.ProductInbound {
	return &ProductService{ob: ob}
}

func (s *ProductService) Save(product *entities.Product) error {
	if err := s.ob.Save(product); err != nil {
		return err
	}
	return nil
}

func (s *ProductService) Find() ([]entities.Product, error) {
	products, err := s.ob.Find()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) FindById(id int) (*entities.Product, error) {
	product, err := s.ob.FindById(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateOne(product *entities.Product, id int) error {
	if err := s.ob.UpdateOne(product, id); err != nil {
		return err
	}
	return nil
}

func (s *ProductService) DeleteOne(id int) error {
	if err := s.ob.DeleteOne(id); err != nil {
		return err
	}
	return nil
}
