package orderAdapter

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	orderPort "github.com/wittawat/go-hex/core/port/order"
	"gorm.io/gorm"
)

// ------------------------ Entities ------------------------
type gormOrder struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	UserID    uint           `gorm:"user_id"`
	ProductID uint           `gorm:"product_id"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Set table name in postgresql database
func (gormOrder) TableName() string {
	return "orders"
}

type gormOrderRepository struct {
	db *gorm.DB
}

// ------------------------ Constructor ------------------------
func NewGormOrderRepository(db *gorm.DB) orderPort.OrderRepository {
	db.AutoMigrate(&gormOrder{})
	if err := orderFactoryPostgres(db); err != nil {
		log.Println("failed to feed order to postgres: ", err)
	}
	return &gormOrderRepository{db: db}
}

// ------------------------ Private Function ------------------------
func getUserFromPostgres(db *gorm.DB) ([]uint, error) {
	var userIDs []uint
	result := db.Model(&entities.User{}).Where("role=?", "user").Pluck("id", &userIDs)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(userIDs) == 0 {
		return nil, errors.New("no user found in users")
	}

	log.Println("order::getUserFromPostgres: success")
	return userIDs, nil
}

func getProductFromPostgres(db *gorm.DB) ([]uint, error) {
	var productIDs []uint
	result := db.Model(&entities.Product{}).Pluck("id", &productIDs)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(productIDs) == 0 {
		return nil, errors.New("no product found in products")
	}

	log.Println("order::getProductFromPostgres: success")
	return productIDs, nil
}

func orderFactoryPostgres(db *gorm.DB) error {
	var count int64
	result := db.Model(&entities.Order{}).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	if count > 0 {
		return nil
	}

	users, err := getUserFromPostgres(db)
	if err != nil {
		return err
	}

	products, err := getProductFromPostgres(db)
	if err != nil {
		return err
	}

	var orders []gormOrder
	for i := 1; i <= 10; i++ {
		randomUser := rand.IntN(len(users))
		randomProd := rand.IntN(len(products))
		order := gormOrder{
			UserID:    users[randomUser],
			ProductID: products[randomProd],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		}
		orders = append(orders, order)
	}
	if result := db.Create(orders); result.Error != nil {
		return result.Error
	}
	log.Println("feed order to postgres: success")
	return nil
}

func entities2GormOrder(o *entities.Order) (*gormOrder, error) {
	userID, err := strconv.Atoi(o.UserID)
	if err != nil {
		return nil, err
	}
	productID, err := strconv.Atoi(o.ProductID)
	if err != nil {
		return nil, err
	}

	var deletedAt gorm.DeletedAt
	if o.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *o.DeletedAt}
	}
	return &gormOrder{
		UserID:    uint(userID),
		ProductID: uint(productID),
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func gorm2EntitiesOrder(gr *gormOrder) *entities.Order {
	id := strconv.FormatUint(uint64(gr.ID), 10)
	userID := strconv.FormatUint(uint64(gr.UserID), 10)
	productID := strconv.FormatUint(uint64(gr.ProductID), 10)
	return &entities.Order{
		ID:        id,
		UserID:    userID,
		ProductID: productID,
		CreatedAt: gr.CreatedAt,
		UpdatedAt: gr.UpdatedAt,
		DeletedAt: &gr.DeletedAt.Time,
	}
}

// ------------------------ Method ------------------------
func (g *gormOrderRepository) Save(ctx context.Context, order *entities.Order) error {
	gr, err := entities2GormOrder(order)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Create(&gr)
	return result.Error
}

func (g *gormOrderRepository) UpdateOne(ctx context.Context, order *entities.Order, id string) error {
	gr, err := entities2GormOrder(order)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Model(&gormOrder{}).Where("id = ?", id).Select("user_id", "product_id", "updated_at").Updates(gr)
	return result.Error
}

func (g *gormOrderRepository) DeleteOne(ctx context.Context, id string) error {
	result := g.db.WithContext(ctx).Delete(&gormOrder{}, id)
	return result.Error
}

func (g *gormOrderRepository) DeleteAllOrderByUser(ctx context.Context, userId string) error {
	userIDStr, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Where("user_id = ?", uint(userIDStr)).Delete(&gormOrder{})
	return result.Error
}

func (g *gormOrderRepository) DeleteAllOrderByProduct(ctx context.Context, productId string) error {
	productIDStr, err := strconv.Atoi(productId)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Where("product_id = ?", uint(productIDStr)).Delete(&gormOrder{})
	return result.Error
}

func (g *gormOrderRepository) FindById(ctx context.Context, id string) (*entities.Order, error) {
	var gr gormOrder
	result := g.db.WithContext(ctx).First(&gr, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return gorm2EntitiesOrder(&gr), nil
}

func (g *gormOrderRepository) FindByUserEmail(ctx context.Context, email string) (*entities.Order, error) {
	var gr gormOrder
	result := g.db.WithContext(ctx).First(&gr).Where("email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return gorm2EntitiesOrder(&gr), nil
}

func (g *gormOrderRepository) FindByUserId(ctx context.Context, userId string) ([]entities.Product, error) {
	var products []entities.Product
	result := g.db.Table("orders").
		Select("products.id, products.title, products.price, products.detail, products.created_at, products.updated_at").
		Joins("JOIN products ON orders.product_id = products.id").
		Where("orders.user_id = ? AND orders.deleted_at IS NULL AND products.deleted_at IS NULL", userId).
		Order("id ASC").
		Scan(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
