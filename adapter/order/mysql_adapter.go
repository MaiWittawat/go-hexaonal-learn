package adapter

import (
	"database/sql"

	"github.com/wittawat/go-hex/core/entities"
)

type MysqlOrderRepository struct {
	db *sql.DB
}

func NewMysqlOrderRepository(db *sql.DB) *MysqlOrderRepository {
	return &MysqlOrderRepository{db: db}
}

func (r *MysqlOrderRepository) Save(order *entities.Order) error {
	query := "INSERT INTO orders (user_id, product_id) VALUES (?, ?)"
	if _, err := r.db.Exec(query, order.UserId, order.ProductId); err != nil {
		return err
	}
	return nil
}

func (r *MysqlOrderRepository) FindByUserId(userId int) ([]entities.Product, error) {
	query := "SELECT p.title, p.price, p.detail FROM orders o JOIN products p ON o.product_id=p.id WHERE o.user_id=?"
	rows, err := r.db.Query(query, userId)
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
	return products, nil
}

func (r *MysqlOrderRepository) UpdateOne(order *entities.Order, id int) error {
	query := "UPDATE orders SET user_id=?, product_id=? WHERE id=?"
	if _, err := r.db.Exec(query, order.UserId, order.ProductId, id); err != nil {
		return err
	}
	return nil
}

func (r *MysqlOrderRepository) DeleteOne(id int) error {
	query := "DELETE FROM orders WHERE id=?"
	if _, err := r.db.Exec(query, id); err != nil {
		return err
	}
	return nil
}
