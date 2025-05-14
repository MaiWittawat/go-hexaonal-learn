package bootstrap

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	jwtAdapter "github.com/wittawat/go-hex/adapter/auth"
	orderAdapter "github.com/wittawat/go-hex/adapter/order"
	productAdapter "github.com/wittawat/go-hex/adapter/product"
	userAdapter "github.com/wittawat/go-hex/adapter/user"
	"github.com/wittawat/go-hex/core/service"
	"github.com/wittawat/go-hex/db"
	"github.com/wittawat/go-hex/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitMongoApp(ctx context.Context, app *gin.Engine) (*mongo.Client, error) {
	mgDB, err := db.InitMongoDB("miniproject", ctx)
	if err != nil {
		log.Fatal("fail to connect mongodb: ", err)
		return nil, err
	}
	log.Println("connected mongo successfully")

	// Authentication
	authNSvc := jwtAdapter.NewAuthNServiceImpl()

	// Repository
	userRepo := userAdapter.NewMongoUserRepository(mgDB.Collection("users"))
	productRepo := productAdapter.NewMongoProductRepository(mgDB.Collection("products"))
	orderRepo := orderAdapter.NewMongoOrderRepository(mgDB.Collection("orders"))

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

	return mgDB.Client(), nil
}
