package orderAdapter

import (
	"context"
	"strconv"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	"gorm.io/gorm"
)

// ------------------------ Entities ------------------------
type GormOrder struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	UserID    uint           `gorm:"user_id"`
	ProductID uint           `gorm:"product_id"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (GormOrder) TableName() string {
	return "orders"
}

type GormOrderRepository struct {
	db *gorm.DB
}

// ------------------------ Constructor ------------------------
func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

// ------------------------ Private Function ------------------------
func entities2GormOrder(o *entities.Order) (*GormOrder, error) {
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
	return &GormOrder{
		UserID:    uint(userID),
		ProductID: uint(productID),
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func gorm2EntitiesOrder(gr *GormOrder) *entities.Order {
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
func (g *GormOrderRepository) Save(ctx context.Context, order *entities.Order) error {
	gr, err := entities2GormOrder(order)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Create(&gr)
	return result.Error
}

func (g *GormOrderRepository) UpdateOne(ctx context.Context, order *entities.Order, id string) error {
	gr, err := entities2GormOrder(order)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Model(&GormOrder{}).Where("id = ?", id).Select("user_id", "product_id", "updated_at").Updates(gr)
	return result.Error
}

func (g *GormOrderRepository) DeleteOne(ctx context.Context, id string) error {
	result := g.db.WithContext(ctx).Delete(&GormOrder{}, id)
	return result.Error
}

func (g *GormOrderRepository) DeleteAllByUser(ctx context.Context, userId string) error {
	userIDStr, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Where("user_id = ?", uint(userIDStr)).Delete(&GormOrder{})
	return result.Error
}

func (g *GormOrderRepository) DeleteAllByProduct(ctx context.Context, productId string) error {
	productIDStr, err := strconv.Atoi(productId)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Where("product_id = ?", uint(productIDStr)).Delete(&GormOrder{})
	return result.Error
}

func (g *GormOrderRepository) FindById(ctx context.Context, id string) (*entities.Order, error) {
	var gr GormOrder
	result := g.db.WithContext(ctx).First(&gr, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return gorm2EntitiesOrder(&gr), nil
}

func (g *GormOrderRepository) FindByUserEmail(ctx context.Context, email string) (*entities.Order, error) {
	var gr GormOrder
	result := g.db.WithContext(ctx).First(&gr).Where("email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return gorm2EntitiesOrder(&gr), nil
}

func (g *GormOrderRepository) FindByUserId(ctx context.Context, userId string) ([]entities.Product, error) {
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
