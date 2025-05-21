package bootstrap

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	jwtAuthNAdapter "github.com/wittawat/go-hex/adapter/auth"
	orderAdapterInbound "github.com/wittawat/go-hex/adapter/order/inbound"
	orderAdapterOutbound "github.com/wittawat/go-hex/adapter/order/outbound"
	productAdapterInbound "github.com/wittawat/go-hex/adapter/product/inbound"
	productAdapterOutbound "github.com/wittawat/go-hex/adapter/product/outbound"
	userAdapterInbound "github.com/wittawat/go-hex/adapter/user/inbound"
	userAdapterOutbound "github.com/wittawat/go-hex/adapter/user/outbound"
	"github.com/wittawat/go-hex/core/services"
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

	redisClient := db.InitRedis()

	// Authentication
	authNSvc := jwtAuthNAdapter.NewAuthenService()

	// Repository
	userRepo := userAdapterOutbound.NewMongoUserRepository(mgDB.Collection("users"))
	productRepo := productAdapterOutbound.NewMongoProductRepository(mgDB.Collection("products"), mgDB.Collection("users"))
	orderRepo := orderAdapterOutbound.NewMongoOrderRepository(mgDB.Collection("orders"), mgDB.Collection("users"), mgDB.Collection("products"))

	// Redis
	userRedisRepo := userAdapterOutbound.NewRedisUserRepository(redisClient, userRepo)
	productRedisRepo := productAdapterOutbound.NewRedisProductRepository(redisClient, productRepo)
	orderRedisRepo := orderAdapterOutbound.NewRedisOrderRepository(redisClient, orderRepo)

	// Authorization
	authZSvc := services.NewAuthZServiceImpl(userRepo)

	// Service by db
	// userService := service.NewUserService(userRepo, productRepo, orderRepo, authNSvc)
	// productService := service.NewProductService(userRepo, productRepo, orderRepo)
	// orderService := service.NewOrderService(orderRepo, userRepo)

	// Service by redis
	userService := services.NewUserService(userRedisRepo, productRedisRepo, orderRedisRepo, authNSvc)
	productService := services.NewProductService(userRedisRepo, productRedisRepo, orderRedisRepo)
	orderService := services.NewOrderService(orderRedisRepo, userRedisRepo)

	// Handler
	userHandler := userAdapterInbound.NewHttpUserHandler(userService)
	productHandler := productAdapterInbound.NewHttpProductHandler(productService)
	orderHandler := orderAdapterInbound.NewHttpOrderHandler(orderService)

	// Routes
	routes.RegisterUserHandler(app, userHandler, authNSvc, authZSvc)
	routes.RegisterProductHandler(app, productHandler, authNSvc, authZSvc)
	routes.RegisterOrderHandler(app, orderHandler, authNSvc, authZSvc)

	return mgDB.Client(), nil
}
