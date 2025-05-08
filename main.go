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

	db, err := database.InitializeDBWithGorm()
	if err != nil {
		log.Fatal("fail to connect mysql with gorm: ", err)

	}
	if err := database.Migration(db); err != nil {
		log.Fatal("fail to migrate database: ", err)
		panic("migrate error")
	}

	app := gin.Default()

	authNSvc := jwtAdapter.NewAuthNServiceImpl()

	userRepo := userAdapter.NewGormUserRepository(db)
	authZSvc := jwtAdapter.NewAuthZServiceImpl(userRepo)

	userService := service.NewUserService(userRepo, authNSvc)
	userHandler := userAdapter.NewHttpUserHandler(userService)
	routes.RegisterUserRoutes(app, userHandler, authNSvc)

	productRepo := productAdapter.NewGormProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := productAdapter.NewHttpProductHandler(productService)
	routes.RegisterProductHandler(app, productHandler, authNSvc, authZSvc)

	orderRepo := orderAdapter.NewGormOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := orderAdapter.NewHttpOrderHandler(orderService)
	routes.RegisterOrderHandler(app, orderHandler)

	if err := app.Run(PORT); err != nil {
		log.Fatal("fail to start server: ", err)
	}
}
