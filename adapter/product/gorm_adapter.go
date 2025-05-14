package productAdapter

import (
	"context"
	"strconv"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	"gorm.io/gorm"
)

// ------------------------ Entities ------------------------
type GormProduct struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Title     string         `gorm:"title"`
	Price     int32          `gorm:"price"`
	Detail    string         `gorm:"detail"`
	CreatedBy uint           `gorm:"created_by"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (GormProduct) TableName() string {
	return "products"
}

type GormProductRepository struct {
	db *gorm.DB
}

// ------------------------ Constructor ------------------------
func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

// ------------------------ Private Function ------------------------
func entities2GormProduct(p *entities.Product) (*GormProduct, error) {
	var deletedAt gorm.DeletedAt
	if p.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *p.DeletedAt}
	}
	userIDint, err := strconv.Atoi(p.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &GormProduct{
		Title:     p.Title,
		Price:     p.Price,
		Detail:    p.Detail,
		CreatedBy: uint(userIDint),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func gorm2EntitiesProduct(gp *GormProduct) *entities.Product {
	id := strconv.FormatUint(uint64(gp.ID), 10)
	userIDStr := strconv.FormatUint(uint64(gp.CreatedBy), 10)
	return &entities.Product{
		ID:        id,
		Title:     gp.Title,
		Price:     gp.Price,
		Detail:    gp.Detail,
		CreatedBy: userIDStr,
		CreatedAt: gp.CreatedAt,
		UpdatedAt: gp.UpdatedAt,
		DeletedAt: &gp.DeletedAt.Time,
	}
}

// ------------------------ Method ------------------------
func (g *GormProductRepository) Save(ctx context.Context, product *entities.Product) error {
	gormProduct, err := entities2GormProduct(product)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Table("products").Create(gormProduct)
	return result.Error
}

func (g *GormProductRepository) Find(ctx context.Context) ([]entities.Product, error) {
	var products []entities.Product
	result := g.db.WithContext(ctx).Table("products").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (g *GormProductRepository) FindById(ctx context.Context, id string) (*entities.Product, error) {
	var gp GormProduct
	result := g.db.WithContext(ctx).Table("products").First(&gp, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return gorm2EntitiesProduct(&gp), nil
}

func (g *GormProductRepository) UpdateOne(ctx context.Context, product *entities.Product, id string) error {
	result := g.db.WithContext(ctx).Table("products").Model(&entities.Product{}).Where("id = ?", id).Select("title", "price", "detail").Updates(product)
	return result.Error
}

func (g *GormProductRepository) DeleteOne(ctx context.Context, id string) error {
	result := g.db.WithContext(ctx).Table("products").Delete(id)
	return result.Error
}

func (g *GormProductRepository) DeleteAll(ctx context.Context, userId string) error {
	userIDInt, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).
		Where("created_by = ?", uint64(userIDInt)).
		Delete(&GormProduct{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
