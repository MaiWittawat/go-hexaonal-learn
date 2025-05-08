package routes

import (
	"github.com/gin-gonic/gin"
	authAdapter "github.com/wittawat/go-hex/adapter/auth"
	middleware "github.com/wittawat/go-hex/adapter/middleware"
	productAdapter "github.com/wittawat/go-hex/adapter/product"
)

func RegisterProductHandler(app *gin.Engine, productHandler *productAdapter.HttpProductHandler, authNSvc *authAdapter.AuthNServiceImpl, authZSvc *authAdapter.AuthZServiceImpl) {
	public := app.Group("/products")
	public.GET("/", productHandler.GetAllProduct)
	public.GET("/:id", productHandler.GetProduct)

	protected := app.Group("/products", middleware.JWTAuthMiddleware(authNSvc))
	protected.POST("/", middleware.RequireRoles(authNSvc, authZSvc, "seller", "admin"), productHandler.CreateProduct)
	protected.PATCH("/:id", middleware.RequireRoles(authNSvc, authZSvc, "seller", "admin"), productHandler.UpdateProduct)
	protected.DELETE("/:id", middleware.RequireRoles(authNSvc, authZSvc, "sellter", "admin"), productHandler.DeleteProduct)
}
