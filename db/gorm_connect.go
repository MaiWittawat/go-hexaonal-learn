package db

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/wittawat/go-hex/core/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDBWithGorm() (*gorm.DB, error) {
	godotenv.Load()
	dns := os.Getenv("MYSQL_URI")
	dialector := mysql.Open(dns)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migration(db *gorm.DB) error {
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entities.Product{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entities.Order{}); err != nil {
		return err
	}
	return nil
}
