package productAdapter

import (
	"context"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	productPort "github.com/wittawat/go-hex/core/port/product"
	"gorm.io/gorm"
)

// ------------------------ Entities ------------------------
type gormProduct struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Title     string         `gorm:"title"`
	Price     int32          `gorm:"price"`
	Detail    string         `gorm:"detail"`
	CreatedBy uint           `gorm:"created_by"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Set table name in postgresql database
func (gormProduct) TableName() string {
	return "products"
}

type gormProductRepository struct {
	db *gorm.DB
}

// ------------------------ Constructor ------------------------
func NewGormProductRepository(db *gorm.DB) productPort.ProductRepository {
	db.AutoMigrate(&gormProduct{})
	// if err := productFactoryPostgres(db); err != nil {
	// 	return nil
	// }
	return &gormProductRepository{db: db}
}

// ------------------------ Private Function ------------------------
func productFactoryPostgres(db *gorm.DB) error {
	var count int64
	result := db.Model(&entities.Product{}).Count(&count)

	if result.Error != nil {
		return result.Error
	}

	if count > 0 {
		return nil
	}

	var products []gormProduct
	for i := 1; i <= 10; i++ {
		iStr := strconv.Itoa(i)
		randomId := rand.IntN(10)
		product := gormProduct{
			Title:     "product" + iStr,
			Price:     int32(i * 10),
			Detail:    "detail for product" + iStr,
			CreatedBy: uint(randomId),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		}
		products = append(products, product)
	}
	db.Create(&products)
	return nil
}

func entities2GormProduct(p *entities.Product) (*gormProduct, error) {
	var deletedAt gorm.DeletedAt
	if p.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *p.DeletedAt}
	}
	userIDint, err := strconv.Atoi(p.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &gormProduct{
		Title:     p.Title,
		Price:     p.Price,
		Detail:    p.Detail,
		CreatedBy: uint(userIDint),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func gorm2EntitiesProduct(gp *gormProduct) *entities.Product {
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
func (g *gormProductRepository) Save(ctx context.Context, product *entities.Product) error {
	gormProduct, err := entities2GormProduct(product)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).Table("products").Create(gormProduct)
	return result.Error
}

func (g *gormProductRepository) Find(ctx context.Context) ([]entities.Product, error) {
	var products []entities.Product
	result := g.db.WithContext(ctx).Table("products").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (g *gormProductRepository) FindById(ctx context.Context, id string) (*entities.Product, error) {
	var gp gormProduct
	result := g.db.WithContext(ctx).Table("products").First(&gp, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return gorm2EntitiesProduct(&gp), nil
}

func (g *gormProductRepository) UpdateOne(ctx context.Context, product *entities.Product, id string) error {
	result := g.db.WithContext(ctx).Table("products").Model(&entities.Product{}).Where("id = ?", id).Select("title", "price", "detail").Updates(product)
	return result.Error
}

func (g *gormProductRepository) DeleteOne(ctx context.Context, id string) error {
	result := g.db.WithContext(ctx).Table("products").Delete(id)
	return result.Error
}

func (g *gormProductRepository) DeleteAll(ctx context.Context, userId string) error {
	userIDInt, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}
	result := g.db.WithContext(ctx).
		Where("created_by = ?", uint64(userIDInt)).
		Delete(&gormProduct{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
