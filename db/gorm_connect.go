package db

import (
	"os"

	"github.com/joho/godotenv"
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
