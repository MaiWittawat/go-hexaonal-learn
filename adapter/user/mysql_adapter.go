package adapter

import (
	"database/sql"

	"github.com/wittawat/go-hex/core/entities"
)

// secondary port
type MysqlUserRepository struct {
	db *sql.DB
}

func NewMysqlUserRepository(db *sql.DB) *MysqlUserRepository {
	return &MysqlUserRepository{db: db}
}

func (r *MysqlUserRepository) Save(user *entities.User) error {
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	return err
}

func (r *MysqlUserRepository) Find() ([]entities.User, error) {
	query := "SELECT username, email, password FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.Username, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, err
}

func (r *MysqlUserRepository) FindById(id int) (*entities.User, error) {
	var user entities.User
	query := "SELECT username, email, password FROM users WHERE id=?"
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MysqlUserRepository) UpdateOne(user *entities.User, id int) error {
	query := "UPDATE users SET username=?, email=?, password=? WHERE id=?"
	if _, err := r.db.Exec(query, user.Username, user.Email, user.Password, id); err != nil {
		return err
	}
	return nil
}

func (r *MysqlUserRepository) DeleteOne(id int) error {
	query := "DELETE FROM users WHERE id=?"
	_, err := r.db.Exec(query, id)
	return err
}
