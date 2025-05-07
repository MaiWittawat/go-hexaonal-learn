package adapter

import (
	"github.com/wittawat/go-hex/core/entities"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(user *entities.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormUserRepository) FindById(id int) (*entities.User, error) {
	var user entities.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) Find() ([]entities.User, error) {
	var users []entities.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *GormUserRepository) UpdateOne(user *entities.User, id int) error {
	result := r.db.Model(&entities.User{}).Where("id = ?", id).
		Select("username", "email", "password").
		Updates(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *GormUserRepository) DeleteOne(id int) error {
	result := r.db.Delete(&entities.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
