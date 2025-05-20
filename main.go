package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/wittawat/go-hex/bootstrap"
	boostrap "github.com/wittawat/go-hex/bootstrap"
)

const PORT = ":3030"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		panic("fail to load env")
	}

	dbType := "postgres" // change dbType to "postgres" if you want to try on postgres
	app := gin.Default()

	switch dbType {
	case "postgres":
		err := boostrap.InitPostgresApp(app)
		if err != nil {
			log.Fatal("failed to init postgres app:", err)
		}
	case "mongo":
		client, err := bootstrap.InitMongoApp(ctx, app)
		if err != nil {
			log.Fatal("failed to init mongo app:", err)
		}
		defer client.Disconnect(ctx)
	default:
		log.Fatal("invalid DB_TYPE")
	}

	if err := app.Run(PORT); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
