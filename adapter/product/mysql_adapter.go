package adapter

import (
	"database/sql"

	"github.com/wittawat/go-hex/core/entities"
)

type MysqlProductRepository struct {
	db *sql.DB
}

func NewMysqlProductRepository(db *sql.DB) *MysqlProductRepository {
	return &MysqlProductRepository{db: db}
}

func (r *MysqlProductRepository) Save(product *entities.Product) error {
	query := "INSERT INTO products (title, price, detail) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, product.Title, product.Price, product.Detail)
	return err
}

func (r *MysqlProductRepository) Find() ([]entities.Product, error) {
	query := "SELECT title, price, detail FROM products"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	var products []entities.Product
	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(&product.Title, &product.Price, &product.Detail); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, err
}

func (r *MysqlProductRepository) FindById(id int) (*entities.Product, error) {
	var product entities.Product
	query := "SELECT title, price, detail FROM products WHERE id=?"
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&product.Title, &product.Price, &product.Detail); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *MysqlProductRepository) UpdateOne(product *entities.Product, id int) error {
	query := "UPDATE products SET title=?, price=?, detail=? WHERE id=?"
	if _, err := r.db.Exec(query, product.Title, product.Price, product.Detail, id); err != nil {
		return err
	}
	return nil
}

func (r *MysqlProductRepository) DeleteOne(id int) error {
	query := "DELETE FROM products WHERE id=?"
	_, err := r.db.Exec(query, id)
	return err
}
