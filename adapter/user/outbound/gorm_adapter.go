// adapter/user/gorm_adapter.go
package userAdapter

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	userPort "github.com/wittawat/go-hex/core/port/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ------------------------ Entities ------------------------
type gormUser struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Role      string         `gorm:"role"`
	Username  string         `gorm:"username"`
	Email     string         `gorm:"email"`
	Password  string         `gorm:"password"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Set table name in postgresql database
func (gormUser) TableName() string {
	return "users"
}

type gormUserRepository struct {
	db *gorm.DB
}

// ------------------------ Constructor ------------------------
func NewGormUserRepository(db *gorm.DB) userPort.UserRepository {
	db.AutoMigrate(&gormUser{})
	// if err := userFactoryPostgres(db); err != nil {
	// 	return nil
	// }
	return &gormUserRepository{db: db}
}

// ------------------------ Private Function ------------------------
func userFactoryPostgres(db *gorm.DB) error {
	var count int64
	result := db.Model(&entities.User{}).Count(&count)

	if result.Error != nil {
		return result.Error
	}

	if count > 0 {
		return nil
	}

	role := []string{"user", "seller", "admin"}
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), 4)
	password := string(hash)

	var users []gormUser
	for i := 1; i <= 10; i++ {
		randomRole := rand.IntN(10) / 4
		fmt.Println("randomRole: ", randomRole)
		user := gormUser{
			Username:  "user" + strconv.Itoa(i),
			Email:     "user@example.com",
			Password:  password,
			Role:      role[randomRole],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		}
		users = append(users, user)
	}
	db.Create(&users)
	return nil
}

func entities2GormUser(u *entities.User) (*gormUser, error) {
	var deletedAt gorm.DeletedAt
	if u.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *u.DeletedAt}
	}
	return &gormUser{
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func gorm2EntitiesUser(gu *gormUser) *entities.User {
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
func (g *gormUserRepository) Save(ctx context.Context, user *entities.User) error {
	gu, err := entities2GormUser(user)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Create(&gu)
	return result.Error
}

func (g *gormUserRepository) UpdateOne(ctx context.Context, user *entities.User, id string) error {
	gu, err := entities2GormUser(user)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Model(&gormUser{}).Where("id = ?", id).Select("username", "email", "password", "updated_at").Updates(gu)
	return result.Error
}

func (g *gormUserRepository) DeleteOne(ctx context.Context, id string) error {
	result := g.db.WithContext(ctx).Delete(&gormUser{}, id)
	return result.Error
}

func (g *gormUserRepository) FindById(ctx context.Context, id string) (*entities.User, error) {
	var gu gormUser
	result := g.db.WithContext(ctx).First(&gu, id)
	if result.Error != nil {
		return nil, result.Error
	}
	user := gorm2EntitiesUser(&gu)
	return user, nil
}

func (g *gormUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var gu gormUser
	result := g.db.WithContext(ctx).Where("email=?", email).First(&gu)
	if result.Error != nil {
		return nil, result.Error
	}
	user := gorm2EntitiesUser(&gu)
	return user, nil
}

func (g *gormUserRepository) Find(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	var gormUsers []gormUser
	result := g.db.WithContext(ctx).Find(&gormUsers)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, gu := range gormUsers {
		users = append(users, *gorm2EntitiesUser(&gu))
	}
	return users, nil
}
