package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	orderAdapter "github.com/wittawat/go-hex/adapter/order"
	productAdapter "github.com/wittawat/go-hex/adapter/product"
	userAdapter "github.com/wittawat/go-hex/adapter/user"
	"github.com/wittawat/go-hex/core/service"
	mysql "github.com/wittawat/go-hex/db"
	"github.com/wittawat/go-hex/routes"
)

const PORT = ":3030"

func main() {

	db, err := mysql.InitializeMysqlDB("mysql")
	if err != nil {
		log.Fatal("fail to connect mysql: ", err)
	}

	app := gin.Default()

	userRepo := userAdapter.NewMysqlUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := userAdapter.NewHttpUserHandler(userService)
	routes.RegisterUserRoutes(app, userHandler)

	productRepo := productAdapter.NewMysqlProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := productAdapter.NewHttpProductHandler(productService)
	routes.RegisterProductHandler(app, productHandler)

	orderRepo := orderAdapter.NewMysqlOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := orderAdapter.NewHttpOrderHandler(orderService)
	routes.RegisterOrderHandler(app, orderHandler)

	if err := app.Run(PORT); err != nil {
		log.Fatal("fail to start server: ", err)
	}
}
