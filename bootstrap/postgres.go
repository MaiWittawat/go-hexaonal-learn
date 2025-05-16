package bootstrap

import (
	"log"

	"github.com/gin-gonic/gin"
	jwtAuthNAdapter "github.com/wittawat/go-hex/adapter/auth"
	orderAdapterInbound "github.com/wittawat/go-hex/adapter/order/inbound"
	orderAdapterOutbound "github.com/wittawat/go-hex/adapter/order/outbound"
	productAdapterInbound "github.com/wittawat/go-hex/adapter/product/inbound"
	productAdapterOutbound "github.com/wittawat/go-hex/adapter/product/outbound"
	userAdapterInbound "github.com/wittawat/go-hex/adapter/user/inbound"
	userAdapterOutbound "github.com/wittawat/go-hex/adapter/user/outbound"
	"github.com/wittawat/go-hex/core/service"
	"github.com/wittawat/go-hex/db"
	"github.com/wittawat/go-hex/routes"
)

func InitPostgresApp(app *gin.Engine) error {
	pgDB, err := db.InitPostgresDB()
	if err != nil {
		return err
	}
	log.Println("connected postgres successfully")

	redisClient := db.InitRedis()

	// Authentication
	authNSvc := jwtAuthNAdapter.NewAuthNServiceImpl()

	// Repository
	userRepo := userAdapterOutbound.NewGormUserRepository(pgDB)
	productRepo := productAdapterOutbound.NewGormProductRepository(pgDB)
	orderRepo := orderAdapterOutbound.NewGormOrderRepository(pgDB)

	// Redis
	userRedisRepo := userAdapterOutbound.NewRedisUserRepository(redisClient, userRepo)
	productRedisRepo := productAdapterOutbound.NewRedisProductRepository(redisClient, productRepo)
	orderRedisRepo := orderAdapterOutbound.NewRedisOrderRepository(redisClient, orderRepo)

	// productRedisRepo := productAdapterOutbound.NewRedisProductRepository(redisClient, productRepo)

	// Authorization
	authZSvc := service.NewAuthZServiceImpl(userRepo)

	// Service
	// userService := service.NewUserService(userRepo, productRepo, orderRepo, authNSvc)
	// productService := service.NewProductService(userRepo, productRepo, orderRepo)
	// orderService := service.NewOrderService(orderRepo, userRepo)

	// Service redis
	userService := service.NewUserService(userRedisRepo, productRedisRepo, orderRedisRepo, authNSvc)
	productService := service.NewProductService(userRedisRepo, productRedisRepo, orderRedisRepo)
	orderService := service.NewOrderService(orderRedisRepo, userRedisRepo)

	// Handler
	userHandler := userAdapterInbound.NewHttpUserHandler(userService)
	productHandler := productAdapterInbound.NewHttpProductHandler(productService)
	orderHandler := orderAdapterInbound.NewHttpOrderHandler(orderService)

	// Routes
	routes.RegisterUserHandler(app, userHandler, authNSvc, authZSvc)
	routes.RegisterProductHandler(app, productHandler, authNSvc, authZSvc)
	routes.RegisterOrderHandler(app, orderHandler, authNSvc, authZSvc)

	return nil
}
