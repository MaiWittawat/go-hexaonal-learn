package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	jwtAdapter "github.com/wittawat/go-hex/adapter/auth"
	orderAdapter "github.com/wittawat/go-hex/adapter/order"
	productAdapter "github.com/wittawat/go-hex/adapter/product"
	userAdapter "github.com/wittawat/go-hex/adapter/user"
	"github.com/wittawat/go-hex/core/service"
	database "github.com/wittawat/go-hex/db"
	"github.com/wittawat/go-hex/routes"
)

const PORT = ":3030"

func main() {

	pgDB, err := database.InitializePgWithGorm()
	if err != nil {
		log.Fatal("fail to connect postgres with gorm: ", err)

	}
	if err := database.Migration(pgDB); err != nil {
		log.Fatal("fail to migrate database: ", err)
	}

	app := gin.Default()

	authNSvc := jwtAdapter.NewAuthNServiceImpl()

	userRepo := userAdapter.NewGormUserRepository(pgDB)
	authZSvc := jwtAdapter.NewAuthZServiceImpl(userRepo)

	userService := service.NewUserService(userRepo, authNSvc)
	userHandler := userAdapter.NewHttpUserHandler(userService)
	routes.RegisterUserRoutes(app, userHandler, authNSvc, authZSvc)

	productRepo := productAdapter.NewGormProductRepository(pgDB)
	productService := service.NewProductService(productRepo)
	productHandler := productAdapter.NewHttpProductHandler(productService)
	routes.RegisterProductHandler(app, productHandler, authNSvc, authZSvc)

	orderRepo := orderAdapter.NewGormOrderRepository(pgDB)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := orderAdapter.NewHttpOrderHandler(orderService)
	routes.RegisterOrderHandler(app, orderHandler, authNSvc, authZSvc)

	if err := app.Run(PORT); err != nil {
		log.Fatal("fail to start server: ", err)
	}
}
