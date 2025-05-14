package bootstrap

import (
	"github.com/gin-gonic/gin"
	jwtAdapter "github.com/wittawat/go-hex/adapter/auth"
	orderAdapter "github.com/wittawat/go-hex/adapter/order"
	productAdapter "github.com/wittawat/go-hex/adapter/product"
	userAdapter "github.com/wittawat/go-hex/adapter/user"
	"github.com/wittawat/go-hex/core/service"
	database "github.com/wittawat/go-hex/db"
	"github.com/wittawat/go-hex/routes"
)

func InitPostgresApp(app *gin.Engine) error {
	pgDB, err := database.InitPostgresDB()
	if err != nil {
		return err
	}
	if err := database.Migration(pgDB); err != nil {
		return err
	}

	// Authentication
	authNSvc := jwtAdapter.NewAuthNServiceImpl()

	// Repository
	userRepo := userAdapter.NewGormUserRepository(pgDB)
	productRepo := productAdapter.NewGormProductRepository(pgDB)
	orderRepo := orderAdapter.NewGormOrderRepository(pgDB)

	// Authorization
	authZSvc := jwtAdapter.NewAuthZServiceImpl(userRepo)

	// Service
	userService := service.NewUserService(userRepo, productRepo, orderRepo, authNSvc)
	productService := service.NewProductService(userRepo, productRepo, orderRepo)
	orderService := service.NewOrderService(orderRepo, userRepo)

	// Handler
	userHandler := userAdapter.NewHttpUserHandler(userService)
	productHandler := productAdapter.NewHttpProductHandler(productService)
	orderHandler := orderAdapter.NewHttpOrderHandler(orderService)

	// Routes
	routes.RegisterUserHandler(app, userHandler, authNSvc, authZSvc)
	routes.RegisterProductHandler(app, productHandler, authNSvc, authZSvc)
	routes.RegisterOrderHandler(app, orderHandler, authNSvc, authZSvc)

	return nil
}
