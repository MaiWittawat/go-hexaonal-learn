package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	orderAdapter "github.com/wittawat/go-hex/adapter/order"
	productAdapter "github.com/wittawat/go-hex/adapter/product"
	userAdapter "github.com/wittawat/go-hex/adapter/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresDB() (*gorm.DB, error) {
	godotenv.Load()
	dns := os.Getenv("POSTGRES_URI")
	log.Println("dns: ", dns)
	dialector := postgres.Open(dns)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migration(db *gorm.DB) error {
	if err := db.AutoMigrate(&userAdapter.GormUser{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&productAdapter.GormProduct{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&orderAdapter.GormOrder{}); err != nil {
		return err
	}
	return nil
}
