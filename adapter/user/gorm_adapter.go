// adapter/user/gorm_adapter.go
package userAdapter

import (
	"context"
	"strconv"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	"gorm.io/gorm"
)

// ------------------------ Entities ------------------------
type GormUser struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Role      string         `gorm:"role"`
	Username  string         `gorm:"username"`
	Email     string         `gorm:"email"`
	Password  string         `gorm:"password"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (GormUser) TableName() string {
	return "users"
}

type GormUserRepository struct {
	db *gorm.DB
}

// ------------------------ Constructor ------------------------
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// ------------------------ Private Function ------------------------
func entities2GormUser(u *entities.User) (*GormUser, error) {
	var deletedAt gorm.DeletedAt
	if u.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *u.DeletedAt}
	}
	return &GormUser{
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func gorm2EntitiesUser(gu *GormUser) *entities.User {
	id := strconv.FormatUint(uint64(gu.ID), 10)
	return &entities.User{
		ID:        id,
		Username:  gu.Username,
		Email:     gu.Email,
		Password:  gu.Password,
		Role:      gu.Role,
		CreatedAt: gu.CreatedAt,
		UpdatedAt: gu.UpdatedAt,
		DeletedAt: &gu.DeletedAt.Time,
	}
}

// ------------------------ Method ------------------------
func (g *GormUserRepository) Save(ctx context.Context, user *entities.User) error {
	gu, err := entities2GormUser(user)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Create(&gu)
	return result.Error
}

func (g *GormUserRepository) UpdateOne(ctx context.Context, user *entities.User, id string) error {
	gu, err := entities2GormUser(user)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Model(&GormUser{}).Where("id = ?", id).Select("username", "email", "password", "updated_at").Updates(gu)
	return result.Error
}

func (g *GormUserRepository) DeleteOne(ctx context.Context, id string) error {
	result := g.db.WithContext(ctx).Delete(&GormUser{}, id)
	return result.Error
}

func (g *GormUserRepository) FindById(ctx context.Context, id string) (*entities.User, error) {
	var gu GormUser
	result := g.db.WithContext(ctx).First(&gu, id)
	if result.Error != nil {
		return nil, result.Error
	}
	user := gorm2EntitiesUser(&gu)
	return user, nil
}

func (g *GormUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var gu GormUser
	result := g.db.WithContext(ctx).Where("email=?", email).First(&gu)
	if result.Error != nil {
		return nil, result.Error
	}
	user := gorm2EntitiesUser(&gu)
	return user, nil
}

func (g *GormUserRepository) Find(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	var gormUsers []GormUser
	result := g.db.WithContext(ctx).Find(&gormUsers)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, gu := range gormUsers {
		users = append(users, *gorm2EntitiesUser(&gu))
	}
	return users, nil
}
